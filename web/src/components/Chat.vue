<template>
  <v-container>
    <v-row class="text-center">
      <v-col class="d-flex" cols="12">
        <v-text-field
          v-model="user"
          label="사용자 명"
        ></v-text-field>
        <v-text-field
          v-model="room"
          label="채팅방 명"
        ></v-text-field>
         <v-btn size="x-large" @click="enter">입장</v-btn>
         <v-btn size="x-large" @click="exit">나가기</v-btn>
      </v-col>
      <v-col class="d-flex" cols="12">
        <v-text-field
          v-model="message"
          label="Message"
        ></v-text-field>
         <v-btn size="x-large" @click="send">전송</v-btn>
      </v-col>
      <v-col cols="12">
        <h2 class="headline font-weight-bold mb-5">
          채팅 내용
        </h2>

         <v-textarea
          v-model="text"
          background-color="amber lighten-4"
          color="orange orange-darken-4"
          label="chat room"
          rows="20"
        ></v-textarea>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>

export default {
  name: 'HelloWorld',

  data: () => ({
    user: "",
    room: "",
    message: "",
    text: "",
    ws: null,
  }),
  methods: {
    enter: function () {
      this.text = "";
      
      if (!this.user || !this.room) {
        this.text += "connect successfully \n";
        return;
      }
      const chatUrl = `ws://localhost:8090/channel?user=${this.user}&room=${this.room}`;
      const ws = new WebSocket(chatUrl);
      const vm = this;
      
      this.ws = ws;
      
      ws.onopen = function () {
        vm.text += "connect successfully \n";
      };

      ws.onclose = function (e) {
        console.log(e)
        vm.text += `closed ${e.code} \n`;
      };

      ws.onmessage = function (e) {
        console.log(e)
        vm.text += `message: ${e.data} \n`;
      };

      ws.onerror = function (e) {
        console.error(e.message);
        vm.text += `error: ${e} \n`;
      };
    },
    exit: function () {
      this.ws.close();
    },
    send: function () {
      try {
        this.ws.send(this.message);
        this.text += `send: ${this.message} \n`;
        this.message = "";
      } catch (e) {
        this.text += `send error: ${e.message} \n`;
      }
    },
  }
}
</script>
