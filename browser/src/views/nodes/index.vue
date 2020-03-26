<template>
  <div class="nodes-container">
    <div
      v-for="i in orderedFullNodes"
      :key="i.node_info.listen_addr"
      :subtitle="getIp(i)"
      class="nodes-item"
    >
      <el-tag type="success">
        <router-link :to="{path:`/Nodes/node/${urlsafeIp(getIp(i))}`}">
          {{ i.node_info.moniker }}
        </router-link>
      </el-tag>
    </div>
  </div>

</template>

<script>
import { mapGetters } from 'vuex'
import { orderBy } from 'lodash'
export default {
  data() {
    return {
      asc: true
    }
  },
  computed: {
    ...mapGetters(['nodes']),
    orderedFullNodes() {
      if (this.nodes) {
        return orderBy(
          this.nodes,
          [n => typeof n.node_info.moniker === 'string' ? n.node_info.moniker.toLowerCase() : n.node_info.moniker],
          this.asc ? 'asc' : 'desc'
        )
      } else {
        return []
      }
    }
  },
  methods: {
    toggleFilter() {
      this.asc = !this.asc
      // this.$store.commit('notify', { title: 'Filtering...', body: 'TODO' })
    },
    toggleSearch() {
      // this.$store.commit('notify', { title: 'Searching...', body: 'TODO' })
    },
    urlsafeIp(ip) {
      return ip && ip.split('.').join('-')
    },
    getIp(fullNode) {
      return (
        fullNode.node_info.listen_addr &&
        fullNode.node_info.listen_addr.split(':')[0]
      )
    }
  }
}
</script>

<style lang="scss" scoped>
.nodes {
  &-container {
    margin-left: 10px;
    margin-right: 30px;
    margin-top: 20px;
    margin-bottom: 20px;
  }
  &-item {
    font-size: 40px;
    line-height: 46px;
  }
}
</style>
