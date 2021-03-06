<template>
  <div id="notebook">
    <aside class="side-bar">
      <div class="toolbar">
        <button @click="addNote" :title="addButtonTitle">
          <i class="material-icons">add</i>
          Add note
        </button>
      </div>
      <div
        class="note"
        v-for="note in notes" :key="note.id"
        @click="selectNote(note)">
        {{ note.title }}
      </div>
    </aside>
    <section class="main">
      <textarea v-model="selectedNote.content"></textarea>
    </section>
    <aside class="preview" v-html="notePreview"></aside>
  </div>
</template>

<script>
  import marked from 'marked'
  export default {
    name: 'App',
    data() {
      return {
        content: "Hello World",
        notes: [],
        selectedId: null,
      }
    },
    computed: {
      notePreview() {
        return this.selectedNote ? marked(this.selectedNote.content) : ''
      },
      addButtonTitle() {
        return `${this.notes.length} note(s) already`
      },
      selectedNote() {
        const note = this.notes.find(
          note => note.id === this.selectedId,
        )
        return note ? note : { content: 'You can write in **markdown**'}
      }
    },
    methods: {
      saveNote() {
        localStorage.setItem('content', this.content)
        this.reportOperation('saving')
      },
      reportOperation(opName) {
        console.log(`The ${opName} operation was completed`)
      },
      selectNote(note) {
        this.selectedId = note.id
      },
      addNote() {
        const time = Date.now()
        const note = {
          id: String(time),
          title: `New note ${this.notes.length + 1}`,
          content: `**Hi!** This note is using [markdown](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet
) for formatting`,
          created: time,
          favorite: false,
        }

        this.notes.push(note)
      }
    },
    watch: {
      content(val) {
        this.saveNote(val)
      }
    },
    created() {
      this.content =
        localStorage.getItem('content')
        || 'You can write in **markdown**'
    }
  }
</script>

<style lang="css">
  .material-icons {
    font-size: 24px;
    line-height: 1;
    vertical-align: middle;
    margin: -3px;
    padding-bottom: 1px;
  }

  #notebook > * {
    float: left;
    display: flex;
    flex-direction: column;
    height: 100%;

    > * {
      flex: auto 0 0;
    }
  }

  .side-bar {
    background: #f8f8f8;
    width: 20%;
    box-sizing: border-box;
  }

  .note {
    padding: 16px;
    cursor: pointer;
  }

  .note:hover {
    background: #ade2ca;
  }

  .note .icon {
    float: right;
  }

  button,
  input,
  textarea {
    font-family: inherit;
    font-size: inherit;
    line-height: inherit;
    box-sizing: border-box;
    outline: none !important;
  }

  button,
  .note.selected {
    background: #40b883;
    color: white;
  }

  button {
    border-radius: 3px;
    border: none;
    display: inline-block;
    padding: 8px 12px;
    cursor: pointer;
  }

  button:hover {
    background: #63c89b;
  }

  input {
    border: solid 2px #ade2ca;
    border-radius: 3px;
    padding: 6px 10px;
    background: #f0f9f5;
    color: #666;
  }

  input:focus {
    border-color: #40b883;
    background: white;
    color: black;
  }

  button,
  input {
    height: 34px;
  }

  .main, .preview {
    width: 40%;
    display: inline-block;
    height: 100%;
  }

  .toolbar {
    padding: 4px;
    box-sizing: border-box;
  }

  .status-bar {
    color: #999;
    font-style: italic;
  }

  textarea {
    resize: none;
    border: none;
    box-sizing: border-box;
    margin: 0 4px;
    font-family: monospace;
  }

  textarea, .notes, .preview {
    flex: auto 1 1;
    overflow: auto;
  }

  .preview {
    padding: 12px;
    box-sizing: border-box;
    border-left: solid 4px #f8f8f8;
  }

  .preview p:first-child {
    margin-top: 0;
  }

  a {
    color: #40b883;
  }

  h1,
  h2,
  h3 {
    margin: 10px 0 4px;
    color: #40b883;
  }

  h1 {
    font-size: 2em;
  }

  h2 {
    font-size: 1.5em;
  }

  h3 {
    font-size: 1.2em;
  }

  h4 {
    font-size: 1.1em;
    font-weight: normal;
  }
</style>
