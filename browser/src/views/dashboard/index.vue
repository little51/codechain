<template>
  <div class="dashboard-container">
    <el-table
      :data="TestnetDate"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="TestnetDate_key"
        label="Testnet Data"
        width="180"
        show-header="false"
      />
      <el-table-column
        prop="TestnetDate_value"
      />
    </el-table>

    <el-table
      :data="CurrentBlock"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="CurrentBlock_key"
        label="Current Block"
        width="180"
      />
      <el-table-column
        prop="CurrentBlock_value"
      />
    </el-table>

    <el-table
      :data="ConnectedTo"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="ConnectedTo_key"
        label="Connected To"
        width="180"
      />
      <el-table-column
        prop="ConnectedTo_value"
      />
    </el-table>
  </div>
</template>

<script>
import num from '../../utils/num'
import { mapGetters } from 'vuex'
import { readableDate } from '../../utils/date'
import votingValidators from '../../utils/votingValidators'

export default {
  computed: {
    ...mapGetters([
      'bc',
      'config',
      'nodes',
      'validators',
      'consensusState',
      'blocks',
      'latestBlock',
      'totalTxs'
    ]),
    validatorsActive() {
      if (this.validators && this.validators.length > 0) {
        return this.validatorCount
      }
      if (this.consensusState && this.consensusState.height_vote_set) {
        return 'STALLED'
      }
      return 'Loading...'
    },
    validatorCount() {
      return `${votingValidators(this.validators).length} voting / ${
        this.validators.length
      } total`
    },
    prevotes() {
      if (this.consensusState && this.consensusState.height_vote_set) {
        let prevotes = this.consensusState.height_vote_set[0].prevotes_bit_array
        let split = prevotes.split(' ')
        let onlineSteak = split[1].split('/')[0]
        let totalSteak = split[1].split('/')[1]
        let minimumSteak = Math.round(totalSteak * 0.6667)
        if (onlineSteak >= minimumSteak) {
          return `${split[3] * 100}% prevoted`
        } else {
          return `${split[3] *
            100}% prevoted (${onlineSteak}steak, need ${minimumSteak}steak)`
        }
      }
      return 'Loading...'
    },
    precommits() {
      if (this.consensusState && this.consensusState.height_vote_set) {
        let precommits = this.consensusState.height_vote_set[0]
          .precommits_bit_array
        let split = precommits.split(' ')
        let onlineSteak = split[1].split('/')[0]
        let totalSteak = split[1].split('/')[1]
        let minimumSteak = Math.round(totalSteak * 0.6667)
        if (onlineSteak >= minimumSteak) {
          return `${Math.round(split[3] * 100)}% precommitted`
        } else {
          return `${Math.round(
            split[3] * 100
          )}% precommitted (${onlineSteak}steak, need ${minimumSteak}steak)`
        }
      }
      return 'Loading...'
    }
  },
  mounted() {
    setTimeout(() => {
      let temp_testnetDate = []
      temp_testnetDate.push({ TestnetDate_key: 'Testnet Version', TestnetDate_value: this.bc.status.node_info.network })
      temp_testnetDate.push({ TestnetDate_key: 'Status', TestnetDate_value: this.validatorsActive })
      temp_testnetDate.push({ TestnetDate_key: 'Prevote State', TestnetDate_value: this.prevotes })
      temp_testnetDate.push({ TestnetDate_key: 'Precommit State', TestnetDate_value: this.precommits })
      temp_testnetDate.push({ TestnetDate_key: 'Total Transactions', TestnetDate_value: num.prettyInt(this.totalTxs) })
      this.TestnetDate = temp_testnetDate

      let temp_current_block = []
      temp_current_block.push({ CurrentBlock_key: 'Block Height', CurrentBlock_value: num.prettyInt(this.latestBlock.height) })
      temp_current_block.push({ CurrentBlock_key: 'Block Time', CurrentBlock_value: this.readableDate(this.latestBlock.time) })
      temp_current_block.push({ CurrentBlock_key: 'Block Height', CurrentBlock_value: num.prettyInt(this.latestBlock.num_txs) })
      temp_current_block.push({ CurrentBlock_key: 'Block Height', CurrentBlock_value: this.latestBlock.last_commit_hash })
      this.CurrentBlock = temp_current_block

      let temp_connectedto = []
      temp_connectedto.push({ ConnectedTo_key: 'RPC Endpoint', ConnectedTo_value: this.bc.rpc })
      this.ConnectedTo = temp_connectedto
    }, 100)
  },
  // eslint-disable-next-line
  data() {
    return {
      num: num,
      TestnetDate: [],
      CurrentBlock: [],
      ConnectedTo: []
    }
  },
  methods: {
    readableDate,
    toggleBlockchainSelect() {
      this.$store.commit(
        'SET_CONFIG_BLOCKCHAIN_SELECT',
        !this.config.blockchainSelect
      )
    }
  }
}
</script>

<style lang="scss" scoped>
.headerCellStyle {
  font-size: 40px;
}
.dashboard {
  &-container {
    margin-left: 10px;
    margin-right: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
