<template>
  <div class="app-container">
    <el-tag v-if="hasPrevPage">
      <router-link :to="{ path: '/Blocks/blocks', query: prevQuery }"><span>&lt; Prev Blocks </span></router-link>
    </el-tag>
    <el-tag v-if="hasNextPage">
      <router-link :to="{ path: '/Blocks/blocks', query: nextQuery }"><span> Next Blocks &gt;</span></router-link>
    </el-tag>

    <el-table
      v-loading="listLoading"
      :data="blocks"
      element-loading-text="Loading"
      fit
      highlight-current-row
    >
      <!-- <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.$index }}
        </template>
      </el-table-column> -->
      <el-table-column label="Height" width="95" :style="{color: '#ECEEF1'}">
        <template slot-scope="scope">
          <el-tag type="info">
            <router-link :to="{path:`/Blocks/blockdetail/${scope.row.header.height}`}">{{ scope.row.header.height }}</router-link>
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Txs" width="50" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.num_txs }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Time" align="center">
        <template slot-scope="scope">
          {{ readableDate(scope.row.header.time) }}
        </template>
      </el-table-column>
      <el-table-column label="Last Commit Hash" align="center">
        <template slot-scope="scope">
          {{ scope.row.header.last_commit_hash }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import axios from 'axios'
import moment from 'moment'
import num from '../../utils/num'

import { mapGetters } from 'vuex'
import { readableDate } from '../../utils/date'

const maxItemsPerPage = 10

export default {
  data() {
    return {
      moment: moment,
      num: num,
      jsonUrl: '',
      itemsPerPage: maxItemsPerPage,
      minHeight: 0,
      maxHeight: 0,
      blocks: [],
      listLoading: true,
      jsonData: null
    }
  },
  computed: {
    ...mapGetters([
      'blockchain',
      'totalBlocks'
    ]),
    hasPrevPage() {
      return this.totalBlocks > 0 && this.minHeight > 1
    },
    hasNextPage() {
      return this.totalBlocks > 0 && this.maxHeight < this.totalBlocks
    },
    prevQuery() {
      if (!this.hasPrevPage) return {}

      let { itemsPerPage, minHeight, maxHeight } = this
      return {
        minHeight: minHeight - itemsPerPage,
        maxHeight: maxHeight - itemsPerPage
      }
    },
    nextQuery() {
      if (!this.hasNextPage) return {}

      let { itemsPerPage, minHeight, maxHeight } = this
      return {
        minHeight: minHeight + itemsPerPage,
        maxHeight: maxHeight + itemsPerPage
      }
    }
  },
  methods: {
    isNext() {
      return this.hasNextPage ? '' : 'info'
    },
    readableDate,
    async fetchPageOfBlocks() {
      let { totalBlocks, itemsPerPage } = this

      let { minHeight, maxHeight } = this.$route.query
      if (minHeight) minHeight = parseInt(minHeight)
      if (maxHeight) maxHeight = parseInt(maxHeight)

      if (totalBlocks > 0 && (!maxHeight || maxHeight < 1 || maxHeight > totalBlocks)) {
        maxHeight = totalBlocks
      }
      if (!minHeight && maxHeight) {
        minHeight = maxHeight - itemsPerPage + 1
      }
      if (minHeight < 1) {
        minHeight = 1
      }

      let query = (minHeight || maxHeight) ? `?minHeight=${minHeight}&maxHeight=${maxHeight}` : ''
      this.jsonUrl = `${this.blockchain.rpc}/blockchain${query}`
      let json = await axios.get(this.jsonUrl)
      this.jsonData = json.data // 所有json数据
      this.blocks = json.data.result.block_metas // json当中的区块数组

      minHeight = 0
      maxHeight = 0
      this.blocks.forEach(block => {
        let height = parseInt(block.header.height)
        if (height < minHeight || !minHeight) minHeight = height
        if (height > maxHeight || !maxHeight) maxHeight = height
      })
      this.minHeight = minHeight
      this.maxHeight = maxHeight
    }
  },
  // eslint-disable-next-line
  async mounted() {
    await this.fetchPageOfBlocks()
    this.listLoading = false
  },
  // eslint-disable-next-line
  watch: {
    '$route'() {
      this.fetchPageOfBlocks()
    }
  }
}
</script>
