<template>
  <main>
    <v-content>
    <v-container grid-list-md text-xs-center>
      <v-layout row wrap>
        <v-flex xs12>          
       <v-card>
          <v-card-text class="blue darken-3">
            <img src="../assets/logo.png">
            <div class="title yellow--text">{{ msg }}<br><br></div>
            <span v-if="error === null">
              <div class="body-1 black-text">
                <b>{{ alias.Fullname }}</b>, the awesome is now activated !</div>
                  <div class="body-1 black--text">
                  Your email address is
                  </div>
                  <div class="title black--text">
                    {{ alias.Alias }}<br><br>
                  </div>
                  <div class="body-1 black--text text-xs-left">
                  You may like to override your name :<br>
                  <v-text-field
                    name="name"    
                    label="name"
                    v-model="alias.Fullname"
                    required></v-text-field>
                  </div>
                  <div class="body-1 black--text text-xs-left">
                  Optionnaly add your PGP public key to enforce email encoding :<br>
            <v-text-field
              name="pgpPubKey"
              label="PGP Pub key"
              v-model="alias.PGPPubKey"
              :rules="pgpRules"                            
              multi-line
            ></v-text-field>
            </div>


            </span>
          </v-card-text>
            <v-card-actions class="blue darken-3" v-if="error === null">
              <v-btn color="yellow" v-on:click="destroyAcc">Delete</v-btn>
              <v-btn color="yellow" v-on:click="updateAcc">Update</v-btn>
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
import VueResource from 'vue-resource'

Vue.use(VueResource)

export default {
  name: 'Validator',
  data () {
    return {
      msg: 'Super !',
      alias: {},
      error: null,
      pgpRules: [
        (v) => !!v || 'Valid PGP public key is required',
        (v) => /-----BEGIN PGP PUBLIC KEY BLOCK-----\n(\w+)*/.test(v) || 'PGP public key must be valid'
      ]
    }
  },
  created: function () {
    this.getData()
  },
  methods: {
    destroyAcc () {
      if (confirm('Do you confirm ?')) {
        let baseUrl = ''
        if (process.env.NODE_ENV === 'development') {
          baseUrl = 'http://localhost:8080'
        }
        var me = this
        Vue.http.delete(baseUrl + '/api/alias/validate/' + this.$route.params.validationKey).then(response => {
          // get body data
          me.entrepriseData = response.body
          window.location = '/'
        }, response => {
          alert('ERR: ' + response.statusText)
        })
      }
    },
    updateAcc () {
      let baseUrl = ''
      if (process.env.NODE_ENV === 'development') {
        baseUrl = 'http://localhost:8080'
      }
      var payload = {
        validationKey: this.$route.params.validationKey,
        fullName: this.alias.Fullname,
        PGPPubKey: btoa(this.alias.PGPPubKey)
      }
      var me = this
      Vue.http.put(baseUrl + '/api/alias/validate/' + this.$route.params.validationKey, payload).then(response => {
        // get body data
        me.entrepriseData = response.body
        // window.location = '/'
      }, response => {
        alert('ERR: ' + response.statusText)
      })
    },
    getData () {
      let baseUrl = ''
      if (process.env.NODE_ENV === 'development') {
        baseUrl = 'http://localhost:8080'
      }
      var me = this
      Vue.http.get(baseUrl + '/api/alias/validate/' + this.$route.params.validationKey).then(response => {
        // get body data
        me.alias = response.body
      }, response => {
        me.msg = response.status + ', ' + response.statusText
        me.error = response.status
      })
    }
  }
}
</script>
