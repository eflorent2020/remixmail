<template>
<div id="app">
  <md-toolbar>

 <h2 class="md-title" style="flex: 1">Admin</h2>

  <md-button class="md-icon-button"  v-if="entrepriseData.LOGGED === 'false'">
    <a v-bind:href="entrepriseData.LOGIN"><md-icon>supervisor_account</md-icon></a>
  </md-button>

  <md-button class="md-icon-button"  v-if="entrepriseData.LOGGED === 'true'">
    <a v-bind:href="entrepriseData.LOGOUT"><md-icon>exit_to_app</md-icon></a>
  </md-button>

  </md-toolbar>
  
  <div class="main-content">
    <h1>{{ msg }}</h1>
  </div>
</div>
</template>

<script>
import Vue from 'vue'

export default {
  name: 'Admin',
  data () {
    return {
      entrepriseData: {},
      msg: '',
      apiKeys: []
    }
  },
  created: function () {
    this.getData()
  },
  methods: {
    getData () {
      let baseUrl = ''
      if (Vue.config.productionTip === false) {
        // baseUrl = 'http://localhost:3000'
      }
      var me = this
      Vue.http.get(baseUrl + '/api/entreprise').then(response => {
        me.entrepriseData = response.body
        if (me.entrepriseData.LOGGED === 'false') {
          this.msg = 'please proceed to login'
        }
      }, response => {
        this.msg = response.status + ' ' + response.statusText
      })
    }
  }
}
</script>
<style>
  .main-content {
  background-color : #111111;
  text-align: center;
  color: white;
  padding-top: 24px;
    background-color: #111111;
  color: #CCC;
  }

body.md-theme-default {
  background-color : #111111;
}



  </style>

