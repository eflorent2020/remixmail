<template>
  <div id="hello">
    <md-toolbar>
    <h1 class="md-title">&nbsp;</h1>
  </md-toolbar>
  
  <div class="main-content">

      <img src="../assets/logo.png">
    <h1>{{ msg }} {{ entrepriseData.APPNAME }}</h1>
    <h2>{{ entrepriseData.TAGLINE }}</h2>
    <ul>
      <li class="register"><h2>Send "REGISTER" to {{ entrepriseData.SERVICE_MAIL }}</h2></li>
  </ul>
  <h2>&nbsp;</h2>
    <ul>
      <li><a href="tos.html">Term of use</a></li>
      <li><a href="privacypolicy.html">Privacy policy</a></li>
      <li><a href="https://github.com/emmanuel-florent/remixmail">Get the source</a></li>
  </ul>
  </div>  
  </div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'

Vue.use(VueResource)

export default {
  name: 'HelloWorld',
  data () {
    return {
      msg: 'Welcome to',
      entrepriseData: {}
    }
  },
  created: function () {
    this.getData()
  },
  methods: {
    getData () {
      let baseUrl = ''
      var me = this
      Vue.http.get(baseUrl + '/api/entreprise').then(response => {
        // get body data
        me.entrepriseData = response.body
        console.log(me.entrepriseData.LOGIN)
      }, response => {
        console.log(response)
        // error callback
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}
.register {
  color: #4CAF50;
}
</style>
