<template>
  <div class="query-container">
    <el-input
      v-model="query_key"
      placeholder="please write key here"
      clearable
      class="query-query"
    />
    <el-button type="primary" @click="queryValue">Query</el-button>
    <textarea
      v-if="json_data!==''"
      v-model="json_data"
      class="query-input"
      readonly
    />
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      query_key: '',
      json_data: '',
      queryData: null
    }
  },
  methods: {
    /**
     * 查询
     */
    async queryValue() {
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
        this.json_data = this.getQueryResponse()
        loading.close()
      }, 500)
    },
    /**
     * 获取文本域的数据
     */
    getQueryResponse() {
      return `\n{\n  "error": "${this.queryData.error}",\n "info": "${this.queryData.info}"\n "result": "${this.queryData.result}"\n}
      `
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
  &-input {
    width: 100%;
    height: 450px;
    border-radius: 15px;
    background-color: #f4f5f7;
    font-weight: 500;
    color: #304156;
    font-family:'Lucida Sans', 'Lucida Sans Regular', 'Lucida Grande', 'Lucida Sans Unicode', Geneva, Verdana, sans-serif;
    border-style: none;
    box-shadow: 1px 1px 5px 1px #888888;
    margin-top: 20px;
    padding-left: 15px;
  }
}
</style>
