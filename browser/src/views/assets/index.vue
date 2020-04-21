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
      <el-form ref="form" label-width="100px">
        <el-form-item label="privatekey:">
          <el-input v-model="form_privatekey" :disabled="signSuccessful" />
        </el-form-item>
        <el-form-item label="publicKey:">
          <el-input v-model="form_publicKey" :disabled="signSuccessful" />
        </el-form-item>
        <el-form-item label="Token:">
          <el-input v-model="form_msg_value_token" :disabled="signSuccessful" />
        </el-form-item>
        <el-form-item label="To:">
          <el-input v-model="form_msg_key_to" :disabled="signSuccessful" />
        </el-form-item>
        <el-form-item label="Amount:">
          <el-input v-model="form_msg_key_amount" :disabled="signSuccessful" />
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
import { Base64 } from 'js-base64'
export default {
  data() {
    return {
      active: 0,
      form_publicKey: '',
      form_privatekey: '',
      signResponse: null,
      registerResponse: null,
      json_data: '',
      register_data: '',
      signSuccessful: false,
      startSign: false,
      startRegister: false,
      resetButton: false,
      form_msg_value_token: '',
      form_msg_key_to: '',
      form_msg_key_amount: ''
    }
  },
  watch: {
    form_msg_key_amount(newData, oldData) {
      if (newData !== '' && !parseInt(newData)) {
        this.open("the information of Amount must be number")
      }
    }
  },
  methods: {
    ResetInfo() {
      this.active = 0
      this.form_publicKey = ''
      this.form_privatekey = ''
      this.signResponse = null
      this.registerResponse = null
      this.json_data = ''
      this.register_data = ''
      this.signSuccessful = false
      this.startSign = false
      this.startRegister = false
      this.resetButton = false
      this.form_msg_value_token = ''
      this.form_msg_key_to = ''
      this.form_msg_key_amount = ''
    },
    /**
     * 构造msg
     */
    getMsgString() {
      let msgValueObj = {
        token: this.form_msg_value_token,
        from: this.form_publicKey,
        to: this.form_msg_key_to,
        amount: this.form_msg_key_amount
      }
      let msgValueString = JSON.stringify(msgValueObj)
      let msgValueBase64String = Base64.encode(msgValueString)
      return msgValueBase64String
    },
    /**
     *  资产登记
     */
    async onRegister() {
      console.log("assets register")
      let postData = {
        "publickey": `${this.form_publicKey}`,
        "sign": `${this.signResponse.sign}`,
        "msg": this.getMsgString()
      }
      let json = await axios.post('http://localhost:3000/assets/new', postData)
      let result = json.data
      this.registerResponse = result
      this.register_data = this.getRegisterResponse()

      // 判断是否正确登记
      if (this.registerResponse.error === '') {
        this.active = 2
      }
      this.resetButton = true
      return Promise.resolve()
    },
    /**
     * 资产签名
     */
    async onSign() {
      console.log('assets sign')
      let postData = {
        "privatekey": `${this.form_privatekey}`,
        "msg": this.getMsgString()
      }
      let json = await axios.post('http://localhost:3000/account/sign', postData)
      let result = json.data
      this.signResponse = result
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
      if (!this.isInputRight()) {
        this.open()
        return
      }
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
     *  判断是否全部正确输入
     */
    isInputRight() {
      if (this.form_privatekey === '' ||
          this.form_publicKey === '' ||
          this.form_msg_value_token === '' ||
          this.form_msg_key_amount === ''
      ) {
        return false
      }
      return true
    },
    /**
     * 未正确输入的提示
     */
    open(params) {
      if (params) {
        this.$message(params)
        return
      }
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
