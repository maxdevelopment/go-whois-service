<template>
  <div id="app">

  </div>
</template>

<script>
  import Routes from '~/static/routes'

  export default {
    name: 'app',
    data: () => ({
      id: null,
      ws: null
    }),
    beforeMount() {
      this.id = Math.random().toString(36).substr(2, 9)
      this.ws = new WebSocket(Routes.socketPath(this.id))
    }
  }
</script>

<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    margin-top: 60px;
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
</style>
