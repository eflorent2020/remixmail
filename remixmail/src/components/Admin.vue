<template>
  <main>
    <v-content>
    <v-container grid-list-md text-xs-center>
      <v-layout row wrap>
        <v-flex xs12>          
       <v-card>
          <v-card-text class="blue darken-3">
            {{ msg }}

<table width="100%">
  <tr v-for="apiKey in apiKeys">
    <td>{{ apiKey.Email }}</td>
    <td>{{ apiKey.ApiKey }}</td>
    <td> <v-icon v-on:click="deleteKey(apiKey)">delete</v-icon></td>     
  </tr>
</table>


          </v-card-text>
            <v-card-actions class="blue darken-3">
              <v-btn color="yellow" v-bind:href="entrepriseData.LOGIN" v-if="entrepriseData.LOGGED === 'false'">Login</Login></v-btn>
              <v-btn color="yellow" v-bind:href="entrepriseData.LOGOUT"  v-if="entrepriseData.LOGGED === 'true'">Logout</v-btn>

              <v-dialog v-model="dialog" persistent max-width="500px">
                    <v-btn color="primary" slot="activator">New API Key</v-btn>
                    <v-card>
                      <v-card-title>
                        <span class="headline">New API Key</span>
                      </v-card-title>
                      <v-card-text>
                        <v-container grid-list-md>
                          <v-layout wrap>
                            <v-flex xs12 sm6 md4>
                              <v-text-field v-model="email"  :rules="emailRules" label="Email address" required></v-text-field>
                            </v-flex>          
                          </v-layout>
                        </v-container>
                        <small>*indicates required field</small>
                      </v-card-text>
                      <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn color="blue darken-1" flat @click.native="dialog = false">Close</v-btn>
                        <v-btn color="blue darken-1" flat  v-on:click="saveForm" @click.native="dialog = false">Save</v-btn>
                      </v-card-actions>
                    </v-card>
              </v-dialog>

            </v-card-actions>
          </v-card>
        </v-flex>
     </v-layout>
    </v-container>
    </v-content>



  </main>
</template>

<script>
import Vue from 'vue'

export default {
  name: 'Admin',
  props: ['entrepriseData'],
  data () {
    return {
      msg: '',
      apiKeys: [],
      dialog: false,
      valid: false,
      email: '',
      emailRules: [
        (v) => !!v || 'E-mail is required',
        (v) => /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'E-mail must be valid'
      ]
    }
  },
  created: function () {
    this.getData()
  },
  methods: {
    deleteKey (apiKey) {
      let baseUrl = '' // myself as root url if empty or ...
      if (process.env.NODE_ENV === 'development') {
        baseUrl = 'http://localhost:8080'
      }
      Vue.http.delete(baseUrl + '/api/keys/' + apiKey.Email).then(response => {
        var me = this
        // let the datastore flush ...
        setTimeout(function () {
          me.getData()
        }, 1000)
      }, response => {
        this.msg = response.status + ' ' + response.statusText
      })
    },
    saveForm () {
      let baseUrl = '' // myself as root url if empty or ...
      if (process.env.NODE_ENV === 'development') {
        baseUrl = 'http://localhost:8080'
      }
      Vue.http.put(baseUrl + '/api/keys/' + this.email).then(response => {
        var me = this
        // let the datastore flush ...
        setTimeout(function () {
          me.getData()
        }, 1000)
      }, response => {
        this.msg = response.status + ' ' + response.statusText
      })
    },
    getData () {
      let baseUrl = '' // myself as root url if empty or ...
      if (process.env.NODE_ENV === 'development') {
        baseUrl = 'http://localhost:8080'
      }
      // var me = this
      Vue.http.get(baseUrl + '/api/keys').then(response => {
        this.apiKeys = response.body
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

td {
  text-align: left;
}

</style>
