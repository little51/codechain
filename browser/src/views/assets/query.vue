<template>
  <div class="query-container">
    <el-input
      v-model="query_key"
      placeholder="please write key here"
      clearable
      class="query-query"
    />
    <el-button type="primary" @click="queryValue">Query</el-button>
    <el-table
      v-if="responseData.length>0"
      border
      :data="responseData"
      style="width: 98%;margin-top:20px"
    >
      <el-table-column
        prop="publickey"
        label="publickey"
        width="900"
      />
      <el-table-column
        prop="token"
        label="token"
        width="180"
      />
      <el-table-column
        prop="amount"
        label="amount"
      />
    </el-table>
    <p v-else align="left" style="color:#97A8BE;padding-left:5px">no Data</p>
  </div>
</template>

<script>
import axios from 'axios'
import { Base64 } from 'js-base64'

export default {
  data() {
    return {
      query_key: '',
      responseData: [],
      queryData: null
    }
  },
  methods: {
    /**
     * 查询
     */
    async queryValue() {
      this.responseData = []
      if (this.query_key === '') {
        this.open()
        return Promise.resolve()
      }
      const loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      let postData = {
        "key": `${this.query_key}`
      }
      setTimeout(async() => {
        let json = await axios.post('http://localhost:3000/assets/query', postData)
        let result = json.data
        this.queryData = result

        // 响应info是字符串，需要json解析
        let resultObj = JSON.parse(result["info"])
        let arrayBase64String = resultObj.result.response.value

        // 拿到array数据
        if (arrayBase64String !== null) {
          let arrayBase64Obj = Base64.decode(arrayBase64String)
          let arrayObj = JSON.parse(arrayBase64Obj)
          let array = arrayObj.array
          this.responseData = array
        }
        loading.close()
      }, 500)
    },
    /**
     * 未正确输入的提示
     */
    open() {
      this.$message('Please complete the information')
    }
  }
}
</script>

<style lang="scss" scoped>
.query {
  &-query {
    width: 450px;
    margin-right: 10px;
  }
  &-container {
    padding: 15px;
  }
}
</style>
