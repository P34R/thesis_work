<template>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">

    <div class=" col-12 row list pt-2 mt-1" >
        <div class="img col-2"><img src="https://media.istockphoto.com/vectors/user-member-vector-icon-for-ui-user-interface-or-profile-face-avatar-vector-id1130884625?k=20&m=1130884625&s=612x612&w=0&h=OITK5Otm_lRj7Cx8mBhm7NtLTEHvp6v3XnZFLZmuB9o=" alt="">
        </div>
        <div class="col-6 " v-on:click="conversation">
           <div class="name m-0">
            <p>{{username}}</p>
           </div>
           <div class="mb-2  ">
           </div>
        </div>
    </div>

</template>

<script>
export default {
  props: ['username', 'pubKey'],
  

  data() {
    return {
      // results: [],
      // userID:null,
    };
  },

   emits: ["startChat"],
  methods:{
    conversation() {
      
        this.$emit('startChat', this.username, this.pubKey);
    },
  },
};


</script>

<style scoped>
.list:hover{
  background: rgba(103, 142, 186, 0.499);

}
    button,
a {
  text-decoration: none;
  padding: 0.2rem 0.5rem;

  font: inherit;
  background-color: #3a0061;
  border: 1px solid #3a0061;
  color: white;
  cursor: pointer;
  border-radius: 30px;
  /* margin-right: 0.5rem; */
  font-size: 10px;
  margin-bottom: 6px;
  /* display: inline-block; */
}

img{
    height: 40px;
    width: 40px;
    border-radius: 50%;

}
.name p{
    /* margin-top: 7px; */
    margin: 0;
  cursor: default;
}
.list{
  margin-top:1px;
  margin-left: 3px;
  background: rgba(153, 142, 186, 0.499);
  padding-top: 4px;
  height: 4rem;
  transition: 0.3s ease-in;
}
</style>