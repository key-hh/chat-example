<template>
  <v-app>
    <v-main>
      <Chat/>
    </v-main>
  </v-app>
</template>

<script>
import Chat from './components/Chat.vue'

export default {
  name: 'App',

  components: {
    Chat,
  },

  data: () => ({
    //
  }),
}
</script>
