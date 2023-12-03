<template>
  <div
    class="mess"
    :class="otherMsg ? 'active' : ''"
  >
    <p style="font-weight: bold; color:black">{{username}} </p>
    <p class="">  {{ message }}</p>
    <p style="text-align: right; font-size: 8px"> {{ stamp.getDate() +'.'+(parseInt(stamp.getMonth())+1)+'.'+stamp.getFullYear()+' '+stamp.getHours()+':'+stamp.getMinutes() }}</p>


<!--    <button v-on:click="deleteChat(Id)" v-if="isHovering">Delete</button> -->
  </div>
</template>

<script>
export default {
  props: ["username", "message", "stamp" ], //chat id

  data() {
    return {
      otherMsg: false,
    };
  },

  methods: {
    setOtherMsg() {
      if (this.username !== this.$store.getters["user/getUsername"]) {
        this.otherMsg = true;
      }

    },
   // func(){},
  },
  mounted() {
    this.setOtherMsg();
  },
};
</script>

<style scoped>
.mess {
  display: inline-block;
  background: rgba(255, 255, 255, 0.719);
  margin-left: 100px;
  padding: 10px;
  /* padding: 20px; */
  border-radius: 10px;
  width:51%;
  margin-top: 10px;
  margin-right:10px;
  margin-bottom:2px;
  float:right;
}
p{
  color:black;
}
.active {
  background: rgba(0, 172, 255, 0.968);

  margin-left:10px;
  float:left;
}
button,
a {
  text-decoration: none;
  padding: 0.3rem 0.5rem;
  font: inherit;
  background-color: #b4302ba1;
  border: 0;
  color: white;
  cursor: pointer;
  border-radius: 30px;
  /* margin-right: 0.5rem; */
  font-size: 10px;
  /* display: inline-block; */
}
</style>
