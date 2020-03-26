<template>
  <div class="node-container">
    <el-table
      :data="ID"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="ID_key"
        label="ID"
        width="180"
        show-header="false"
      />
      <el-table-column
        prop="ID_value"
      />
    </el-table>

    <el-table
      :data="PubKey"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="PubKey_key"
        label="Pub Key"
        width="180"
        show-header="false"
      />
      <el-table-column
        prop="PubKey_value"
      />
    </el-table>

    <el-table
      :data="NetWork"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="NetWork_key"
        label="NetWork"
        width="180"
        show-header="false"
      />
      <el-table-column
        prop="NetWork_value"
      />
    </el-table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { readableDate } from '../../utils/date'

export default {
  data: () => ({
    tmpFullNode: {
      node_info: {
        moniker: 'Loading...',
        pub_key: 'todoreplacemewithfullNodepubkey'
      }
    },
    persistentPeer: null,
    ID: [],
    PubKey: [],
    NetWork: []
  }),
  computed: {
    ...mapGetters(['nodes']),
    fullNode() {
      if (this.nodes && this.nodes.length > 0) {
        return this.nodes.find(
          v =>
            v.node_info.listen_addr ===
            this.$route.params.node + '://0.0.0.0:26656'
        )
      } else {
        return this.tmpFullNode
      }
    }
  },

  mounted() {
    this.persistentPeer = `${this.fullNode.node_info.id}@${this.fullNode.node_info.listen_addr}`

    this.ID.push({ ID_key: 'Moniker', ID_value: this.fullNode.node_info.moniker })
    this.ID.push({ ID_key: 'IP', ID_value: this.getIp(this.fullNode) })
    this.ID.push({ ID_key: 'Start Date', ID_value: this.fullNode.connection_status && readableDate(this.fullNode.connection_status.SendMonitor.Start) })

    this.PubKey.push({ PubKey_key: 'Value', PubKey_value: this.fullNode.node_info.id })
    this.PubKey.push({ PubKey_key: 'Persistent Peer', PubKey_value: this.persistentPeer })

    this.NetWork.push({ NetWork_key: 'Network', NetWork_value: this.fullNode.node_info.network })
    this.NetWork.push({ NetWork_key: 'Tendermint Version', NetWork_value: this.fullNode.node_info.version })
    this.NetWork.push({ NetWork_key: 'Channels', NetWork_value: this.fullNode.node_info.channels })
  },
  methods: {
    urlsafeIp(ip) {
      return ip.split('.').join('-')
    },
    getIp(fullNode) {
      return fullNode.node_info.listen_addr
    },
    readableDate
  }
}
</script>

<style lang="scss" scoped>
.node {
  &-container {
    margin-left: 10px;
    margin-right: 30px;
  }
}
</style>
