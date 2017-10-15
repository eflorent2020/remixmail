<template>
  
  <span class="hello">
    <md-toolbar>
    <h1 class="md-title">&nbsp;</h1>
  </md-toolbar>
      <center><img src="../assets/logo.png"></center>
      <h1>{{ msg }}</h1>
    <span v-if="error === null">
    <h1>Congrats <span class="name">{{ alias.Fullname }}</span> everyhing is now activated</h1>
    <h2> your mail address is</h2>    
    <h2 class="alias">{{ alias.Alias }}</h2>
    <ul>
      <li>You may like to override your name 
        <input type="text" name="name" v-model="alias.Fullname">
                <input  v-on:click="updateAcc" type="submit">
        </li>       
  </ul>
      <h2>&nbsp;</h2> 
    <p>You may like to delete your account <br>
        <input type="button" v-on:click="destroyAcc" class="delete" value="delete">                
        </p>       
        </span>
  </span>
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
      error: null
    }
  },
  created: function () {
    this.getData()
  },
  methods: {
    destroyAcc () {
      if (confirm('Do you confirm ?')) {
        let baseUrl = ''
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
      var payload = {
        validationKey: this.$route.params.validationKey,
        fullName: this.alias.Fullname
      }
      var me = this
      Vue.http.put(baseUrl + '/api/alias/validate/' + this.$route.params.validationKey, payload).then(response => {
        // get body data
        me.entrepriseData = response.body
        window.location = '/'
      }, response => {
        alert('ERR: ' + response.statusText)
      })
    },
    getData () {
      let baseUrl = ''
      var me = this
      Vue.http.get(baseUrl + '/api/alias/validate/' + this.$route.params.validationKey).then(response => {
        // get body data
        me.alias = response.body
      }, response => {
        me.msg = response.status + ' ' + response.statusText
        me.error = response.status
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
body.md-theme-default {
  background-color : #111111;
}
.hello {
  background-color : #111111;
  text-align: center;
  color: white;
  padding-top: 24px;
    background-color: #111111;
  color: #CCC;
}
h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}


.alias {
  color: #F44336;
  font-weight: bold;
  font-size: 18px;
}

input[type="submit"] {
  background-color: green;
}

input[type="button"] {
  padding: 10px;
  margin-top: 12px;
  border: solid 1px #dcdcdc;
  transition: box-shadow 0.3s, border 0.3s;
  background-color: red;
  color: black;  
}

input[type="text"], input[type="submit"] {
  padding: 10px;
  border: solid 1px #dcdcdc;
  transition: box-shadow 0.3s, border 0.3s;
  color: black;
}
input[type="text"]:focus,
input[type="text"].focus {
  border: solid 1px #707070;
  box-shadow: 0 0 5px 1px #969696;

}

.name {
  color: white;
}

</style>
