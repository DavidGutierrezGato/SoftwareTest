<script setup>
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import VueSplitView from 'vue-split-view'

</script>

<script>
export default {
  data() {
    return {
      searchValue: '',
      data2: [
        { id: 1, column1: 'Value 1', column2: 'Value 2', column3: 'Value 3' },
        { id: 2, column1: 'Value 4', column2: 'Value 5', column3: 'Value 6' }
      ],
      data3:[],
      selectedItem:null,
      selectedRow: null,
      page:0
    }
  },
  computed: {
    uniqueData() {
      return this.data3 = new Set(this.data3)
      //this.data3.;
    },
  
  },
  methods: {
    find() {
      // Perform the search with the value of searchValue
      
      console.log(`Searching for ${this.searchValue}`)

      fetch('http://localhost:3000/find/'+this.searchValue+'/'+this.page)
      .then(response => response.json())
      .then(data3 => this.data3 = data3)
      
      this.selectedItem = null;
     
    },
    select(item,index) {
      this.selectedItem = item.source.Body;
      this.selectedRow = index;
    },
    addPage(){
      this.page += 20
      this.selectedItem = null;
      this.selectedRow = null;
      this.find()
    },
    subPage(){
      this.page -= 20
      if(this.page < 0)
      {
        this.page = 0
      }
      this.selectedItem = null;
      this.selectedRow = null;
      this.find()
    }

    
    
  }
}


  


</script>




<template>
  


 
 
 


  <body>

    <div class="container">
      <input type="text" v-model="searchValue" placeholder="Enter search term">
     <button v-on:click="find()">Find</button>
    
    </div>
    
    
    <br>
    <br>

    <div class="container">

      <table>
    <thead>
      <tr>
        <th>From</th>
        <th>To</th>
        <th>Subject</th>
        
      </tr>
    </thead>
    <tbody>
      <tr v-for="(item,index) in uniqueData" :key="item.id"  v-bind:class="{ 'selected': item === selectedItem, 'selected':selectedRow === index }" @click="select(item,index)" 
      >
        <td>{{ item.source.From}}</td>
        <td>{{ item.source.To}}</td>
        <td>{{ item.source.Subject}}</td>
        
        
      </tr>
    </tbody>
  </table>
  <button v-on:click="subPage()">Prev</button>
  <button v-on:click="addPage()">Next</button>

    </div>
    
    
  
    <div class="container">
  
  <p class="custom-paragraph">
      {{ this.selectedItem}}
    </p>

</div>
  
  
  

  </body>
  

  
</template>

<style >



.custom-paragraph {
  font-size: 18px;
  font-weight: 500;
  color: rgb(0, 0, 0);
  text-align: justify;
  
}

input[type='text'] {
  padding: 12px 20px;
  margin: 8px 0;
  box-sizing: border-box;
  border: 2px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

button {
  width: 100%;
  background-color: #4CAF50;
  color: white;
  padding: 14px 20px;
  margin: 8px 0;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

button:hover {
  background-color: #45a049;
}

.container {
    display: flex;
    flex-direction: column;
}



table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  border: 1px solid black;
  padding: 8px;
  text-align: left;
}

th {
  background-color: #f2f2f2;
}


.selected {
  background-color: #06e79c;
}



</style>
