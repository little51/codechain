<template>
  <div class="assets-container">
    <div class="assets-title">
      <div style="width: 400px">
        <el-steps :active="active" finish-status="success">
          <el-step title="Asset signature" icon="el-icon-edit" />
          <el-step title="Add New Account" icon="el-icon-upload" />
        </el-steps>
      </div>
    </div>
    <div>
      <el-form ref="form" label-width="80px">
        <el-form-item label="privatekey">
          <el-input v-model="form_privatekey" :disabled="signSuccessful" />
        </el-form-item>
        <el-form-item label="publicKey">
          <el-input v-model="form_publicKey" :disabled="signSuccessful" />
        </el-form-item>
        <el-form-item label="msg_key">
          <el-input v-model="form_msg_key" :disabled="signSuccessful" />
        </el-form-item>
        <el-form-item label="msg_value">
          <el-input v-model="form_msg_value" :disabled="signSuccessful" />
        </el-form-item>
        <el-form-item v-show="!signSuccessful">
          <el-button type="primary" @click="onSubmit">Asset Register</el-button>
        </el-form-item>
        <el-button v-if="resetButton" type="primary" class="assets-reset" @click="ResetInfo">Reset</el-button>
      </el-form>
      <p v-show="startSign" class="assets-showResponse">The Response of Asset signature</p>
      <textarea
        v-show="startSign"
        v-model="json_data"
        class="assets-input1"
        readonly
      />
      <p v-show="startRegister" class="assets-showResponse">The Response of Add New Account</p>
      <textarea
        v-show="startRegister"
        v-model="register_data"
        class="assets-input2"
        readonly
      />
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      active: 0,
      form_publicKey: '',
      form_privatekey: '',
      form_msg_key: '',
      form_msg_value: '',
      signResponse: null,
      registerResponse: null,
      json_data: '',
      register_data: '',
      signSuccessful: false,
      startSign: false,
      startRegister: false,
      resetButton: false
    }
  },
  methods: {
    ResetInfo() {
      this.active = 0
      this.form_publicKey = ''
      this.form_privatekey = ''
      this.form_msg_key = ''
      this.form_msg_value = ''
      this.signResponse = null
      this.registerResponse = null
      this.json_data = ''
      this.register_data = ''
      this.signSuccessful = false
      this.startSign = false
      this.startRegister = false
      this.resetButton = false
    },
    /**
     *  资产登记
     */
    async onRegister() {
      if (this.form_publicKey === '' || this.form_privatekey === '' || this.register_msg_key === '' || this.register_msg_value === '') {
        this.open()
        return Promise.resolve()
      }
      let postData = {
        "publickey": `${this.form_publicKey}`,
        "sign": `${this.signResponse.sign}`,
        "msg": `{\"key\": \"${this.form_msg_key}\",\"value\": \"${this.form_msg_value}\"}`
      }
      let json = await axios.post('http://localhost:3000//assets/new', postData)
      let result = json.data
      this.registerResponse = result
      console.log(this.registerResponse)
      this.register_data = this.getRegisterResponse()

      // 判断是否正确请求成功
      if (this.registerResponse.error === '') {
        this.resetButton = true
        this.active = 2
      }
      return Promise.resolve()
    },
    /**
     * 资产签名
     */
    async onSign() {
      if (this.form_privatekey === '' || this.form_msg_key === '' || this.form_msg_value === '' || this.form_publicKey === '') {
        this.open()
        return Promise.resolve()
      }
      let postData = {
        "privatekey": `${this.form_privatekey}`,
        "msg": `{\"key\": \"${this.form_msg_key}\",\"value\": \"${this.form_msg_value}\"}`
      }
      let json = await axios.post('http://localhost:3000/account/sign', postData)
      let result = json.data
      this.signResponse = result
      console.log(this.signResponse)
      this.json_data = this.getSignResponse()

      // 判断是否正确请求成功
      if (this.signResponse.error === '') {
        this.signSuccessful = true
        this.active = 1
      }
      return Promise.resolve()
    },
    /**
     * 提交按钮
     */
    async onSubmit() {
      const loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      setTimeout(async() => {
        this.startSign = true
        await this.onSign()
        this.startRegister = true
        await this.onRegister()
        loading.close()
      }, 1000)
    },
    /**
     * 未正确输入的提示
     */
    open() {
      this.$message('Please complete the information')
    },
    /**
     * 获取资产签名的响应结果的json
     */
    getSignResponse() {
      return `\n{\n  "error": "${this.signResponse.error}",\n  "sign": "${this.signResponse.sign}"\n}
      `
    },
    /**
     * 获取资产登记的响应结果的json
     */
    getRegisterResponse() {
      return `\n{\n  "error": "${this.registerResponse.error}",\n  "info": "${this.registerResponse.info}",\n  "result": ${this.registerResponse.result}\n}`
    }
  }
}
</script>

<style lang="scss" scoped>
.assets {
  &-container {
    margin-left: 10px;
    margin-right: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
  &-title {
    margin-left: 10px;
    margin-top: 20px;
    margin-bottom: 30px;
  }
  &-input1 {
    width: 100%;
    height: 120px;
    margin-top: 1px;
    margin-bottom: 10px;
    border-radius: 5px;
    background-color: #f4f5f7;
    font-weight: 500;
    padding-left: 20px;
    color: #304156;
    font-family:'Lucida Sans', 'Lucida Sans Regular', 'Lucida Grande', 'Lucida Sans Unicode', Geneva, Verdana, sans-serif;
    border-style: none;
    box-shadow: 1px 1px 5px 1px #888888
  }
  &-input2 {
    width: 100%;
    height: 620px;
    margin-top: 1px;
    margin-bottom: 10px;
    border-radius: 5px;
    background-color: #f4f5f7;
    font-weight: 500;
    padding-left: 20px;
    color: #304156;
    font-family:'Lucida Sans', 'Lucida Sans Regular', 'Lucida Grande', 'Lucida Sans Unicode', Geneva, Verdana, sans-serif;
    border-style: none;
    box-shadow: 1px 1px 5px 1px #888888
  }
  &-next {
    margin-top: 20px;
    margin-left: 95%;
  }
  &-reset {
    margin-left: 95%;
  }
  &-register {
    width: 100%;
    height: 600px;
    border: none
  }
  &-showResponse {
    font-weight: 700;
    color: #908989;
    font-size: 16px;
    line-height: 38px;
  }
}
</style>
