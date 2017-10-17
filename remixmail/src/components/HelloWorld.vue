<template>
  <main>
    <v-content>
    <v-container grid-list-md text-xs-center>
      <v-layout row wrap>
        <v-flex xs12>          
       <v-card>
          <v-card-text class="blue darken-3">
        <img src="../assets/logo.png">
        <div class="title yellow--text">{{ msg }}<br>{{ entrepriseData.APPNAME }}</div>
       <div class="sub-heading black--text"><i>{{ entrepriseData.TAGLINE }}</i></div>
       -
          <div class="body-1 black--text">Write to<br>
            <a v-bind:href="'mailto:' + entrepriseData.SERVICE_MAIL + '?Subject=register'">
             
                  <v-btn color="primary"> {{ entrepriseData.SERVICE_MAIL }}</v-btn>
              </a><br> 
            with subject : "REGISTER" to activate the awesome
            </div>
-

          </v-card-text>
</v-card>


        </v-flex>
     </v-layout>
    </v-container>
    </v-content>
  </main>
</template>

<script>

export default {
  props: ['entrepriseData'],
  data () {
    return {
      msg: 'Welcome to'
    }
  }
}
</script>
