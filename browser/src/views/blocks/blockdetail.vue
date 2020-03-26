<template>
  <div class="app-container">
    <el-table
      :data="BlockHeader"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="BlockHeader_key"
        label="Block Header"
        width="180"
        show-header="false"
      />
      <el-table-column
        prop="BlockHeader_value"
      />
    </el-table>

    <el-table
      :data="LastBlock"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="LastBlock_key"
        label="Last Block"
        width="180"
        show-header="false"
      />
      <el-table-column
        prop="LastBlock_value"
      />
    </el-table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import axios from 'axios'
import createHash from 'create-hash'
import varint from 'varint'
import b64 from 'base64-js'
import { decodeBase64, decodeTx } from '../../utils/tx'

export default {
  data: () => ({
    BlockHeader: [],
    LastBlock: [],
    jsonUrl: '',
    block: {
      header: {
        chain_id: '',
        height: 0,
        time: '',
        num_txs: 0,
        last_block_id: {
          hash: '',
          parts: {
            total: 0,
            hash: ''
          }
        },
        last_commit_hash: '',
        data_hash: '',
        validators_hash: '',
        app_hash: ''
      },
      data: {
        txs: []
      },
      last_commit: {
        blockID: {
          hash: '',
          parts: {
            total: 0,
            hash: ''
          }
        },
        precommits: [
          {
            validator_address: '',
            validator_index: 0,
            height: 0,
            round: 0,
            type: 0,
            block_id: {
              hash: '',
              parts: {
                total: 0,
                hash: ''
              }
            },
            signature: [0, '']
          }
        ]
      }
    }
  }),
  computed: {
    ...mapGetters([
      'blockchain',
      'totalBlocks',
      'latestBlock'
    ]),
    prevHeight() {
      return this.block.header.height - 1
    },
    nextHeight() {
      return this.block.header.height + 1
    },
    hasPrevBlock() {
      return this.prevHeight > 0
    },
    hasNextBlock() {
      return this.nextHeight <= this.totalBlocks
    },
    decodedTxs() {
      return this.block.data.txs.map((tx, i) => {
        let txObj = decodeTx(tx)
        let hash = this.txHash(i)
        let txHash = {
          isRouterLink: true,
          title: 'View transaction details',
          text: hash,
          to: { name: 'tx', params: { hash }}
        }
        return Object.assign({ txHash }, txObj)
      })
    }
  },
  watch: {
    // eslint-disable-next-line
    '$route'(to, from) {
      this.fetchBlock()
    }
  },
  async mounted() {
    await this.fetchBlock()
    this.BlockHeader.push({ BlockHeader_key: 'Chain ID', BlockHeader_value: this.block.header.chain_id })
    this.BlockHeader.push({ BlockHeader_key: 'Time', BlockHeader_value: this.block.header.time })
    this.BlockHeader.push({ BlockHeader_key: 'Transactions', BlockHeader_value: this.block.header.num_txs })
    this.BlockHeader.push({ BlockHeader_key: 'Last Commit Hash', BlockHeader_value: this.block.header.last_commit_hash })
    this.BlockHeader.push({ BlockHeader_key: 'Validators Hash', BlockHeader_value: this.block.header.validators_hash })
    this.BlockHeader.push({ BlockHeader_key: 'App Hash', BlockHeader_value: this.block.header.app_hash })

    this.LastBlock.push({ LastBlock_key: 'Hash', LastBlock_value: this.block.header.last_block_id.hash })
    this.LastBlock.push({ LastBlock_key: 'Parts Total', LastBlock_value: this.block.header.last_block_id.parts.total })
    this.LastBlock.push({ LastBlock_key: 'Parts Hash', LastBlock_value: this.block.header.last_block_id.parts.hash })
  },
  methods: {
    txHash(idx) {
      let tx = this.block.data.txs[idx]
      let b64str = tx.replace(/^:base64:/, '')
      let buffer = Buffer.from(b64str, 'base64')
      let hex = createHash('sha256').update(buffer).digest('hex')
      return hex.substr(0, 40).toUpperCase()
    },
    async fetchBlock() {
      if (this.$route.params.block) {
        this.jsonUrl = `${this.blockchain.rpc}/block?height=${this.$route.params.block}`
      } else {
        this.jsonUrl = `${this.blockchain.rpc}/block?height=${this.latestBlock.height}`
      }
      // console.log(this.jsonUrl)
      let json = await axios.get(this.jsonUrl)
      this.block = json.data.result.block
      this.block.header.height = parseInt(this.block.header.height)
    },

    // TODO deprecate? (yes)
    queryTxs() {
      return this.queryTx(this.block.data.txs.length)
    },

    // TODO deprecate? (yes)
    queryTx(len, key = 0) {
      return new Promise((resolve, reject) => {
        if (key >= len) resolve()
        let txstring = decodeBase64(this.block.data.txs[key])
        // console.log(txstring)
        let txbytes = b64.toByteArray(this.block.data.txs[key])
        // console.log(txbytes)
        let varintlen = new Uint8Array(varint.encode(txbytes.length))
        // console.log(varintlen)
        let tmp = new Uint8Array(varintlen.byteLength + txbytes.byteLength)
        tmp.set(new Uint8Array(varintlen), 0)
        tmp.set(new Uint8Array(txbytes), varintlen.byteLength)
        // console.log(tmp)

        // TODO replace with fixed decoding that uses SHA256
        let hash = createHash('ripemd160')
          .update(Buffer.from(tmp))
          .digest('hex')
        // console.log(hash)

        let url = `${this.blockchain.rpc}/tx?hash=0x${hash}`
        axios
          .get(url)
          .then(json => {
            // console.log(json)
            json.data.result.string = txstring
            this.block.data.txs.splice(key, 1, json.data.result)
            this.queryTx(len, key + 1)
              .then(resolve)
              .catch(reject)
          })
          .catch(reject)
      })
    }
  }
}
</script>
