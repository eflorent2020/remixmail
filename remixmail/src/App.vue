<template>
  <div id="app" dark>
    <v-app toolbar footer>
    <v-toolbar class="blue darken-4">
      <v-toolbar-title v-text=""></v-toolbar-title>
    </v-toolbar>
    <router-view v-bind:entrepriseData="entrepriseData"/>

  <v-footer app class="blue darken-4" dark>
          <v-layout row wrap  >
            <v-flex xs12 align-center text-xs-center class="white--text">
               ~ 
                <a class="white--text" href="tos.html">Term of use</a> ~ 
                <a class="white--text" href="privacypolicy.html">Privacy policy</a> ~ 
                <a class="white--text" href="https://github.com/emmanuel-florent/remixmail">Get the source</a> 
               ~ 
            </v-flex>
          </v-layout>
  </v-footer>
    </v-app>
  </div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'

Vue.use(VueResource)

export default {
  name: 'app',
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
      let baseUrl = '' // myself as root url if empty or ...
      if (process.env.NODE_ENV === 'development') {
        baseUrl = 'http://localhost:8080'
      }
      var me = this
      Vue.http.get(baseUrl + '/api/entreprise').then(response => {
        // get body data
        me.entrepriseData = response.body
      }, response => {
        // error callback
        console.log(response)
      })
    }
  }
}
</script>
<style>
.content {
  background: rgb(25,118,210);
}

.my-card {
  background: rgb(16, 58, 106);
}

</style>
<style lang="stylus">
  @require './main'
</style>
