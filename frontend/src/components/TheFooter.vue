<template>
  <footer>
    <span class="todo-count">{{ count }} items left</span>
    <ul>
      <li>
        <button
          @click="filterChange('all')"
          id="all"
          :class="clicked == 'all' ? 'clicked' : ''"
        >
          All
        </button>
      </li>
      <li>
        <button
          @click="filterChange('active')"
          id="active"
          :class="clicked == 'active' ? 'clicked' : ''"
        >
          Active
        </button>
      </li>
      <li>
        <button
          @click="filterChange('completed')"
          id="completed"
          :class="clicked == 'completed' ? 'clicked' : ''"
        >
          Completed
        </button>
      </li>
    </ul>
  </footer>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  emits: ["change-filter"],
  data() {
    return {
      clicked: "all" as string,
    };
  },
  props: {
    count: {
      type: Number,
      required: true,
      default: 0,
    },
  },
  methods: {
    filterChange(clickedVal: string) {
      this.clicked = clickedVal;
      this.$emit("change-filter", this.clicked);
    },
  },
});
</script>

<style scoped>
footer {
  color: #777;
  padding: 1rem;
  text-align: center;
  border-top: 1px solid #e6e6e6;
}
.todo-count {
  float: left;
  text-align: left;
}
ul {
  margin: 0;
  padding: 0;
  list-style: none;
}
li {
  display: inline;
  padding: 1rem;
}
button {
  color: rgb(143, 142, 142);
  text-decoration: none;
  border: none;
  background-color: transparent;
}
button:hover {
  cursor: pointer;
  color: black;
}
.clicked {
  color: black;
}
</style>
