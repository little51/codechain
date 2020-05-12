<template>
  <div>
    <input v-model="privatekey" type="text">
    <input v-model="publickey" type="text">
    <input v-model="msg" type="text">
    <button @click="handleClick">commit</button>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
      privatekey: "",
      publickey: "",
      msg: ""
    }
  },
  methods: {
    async handleClick() {
      let postData = {
        "privatekey": this.privatekey,
        "publickey": this.publickey,
        "msg": this.msg
      }
      let json = await axios.post('http://localhost:4000/broadcast/msg', postData)
      console.log(json)
    }
  }
}
</script>
