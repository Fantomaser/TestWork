var app = new Vue({
    el: '#app',
    data: {
      message: null,
      btnName: 'Сгенерировать',

      userInfo: null,

      OutpKeyMess: '0',
      OutpURLMess: '/' 
    },
    methods:{
      CheckText: function(){
        console.log('Chekau', this.message)
              // GET /someUrl
        this.$http.post('http://localhost:8080/makeText', {msg : this.message}).then(res => {
            console.log('Generation message', res)
            this.userInfo = res.body.user
            this.OutpKeyMess = this.userInfo.Key
            this.OutpURLMess = this.userInfo.Adres
        }, res => {
            console.log("SendError!!!!!!", res)
        });
      }
    }
  })