<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'

// Types
interface WebSocketMessage {
  type: 'edit'
  data: string
  sender_id?: string
}

type ConnectionStatus = 'Connected' | 'Connecting' | 'Disconnected'

// State
const documentText = ref<string>("")
const status = ref<ConnectionStatus>("Disconnected")
const socket = ref<WebSocket | null>(null)

// Connect to WebSocket
onMounted(() => {
  status.value = "Connecting"
  
  // Connect to your Go Backend
  socket.value = new WebSocket("ws://localhost:8080/ws")
  
  // Connection opened
  socket.value.onopen = () => {
    status.value = "Connected"
    console.log("Connected to SyncPad Server")
  }
  
  // Listen for messages (From other users)
  socket.value.onmessage = (event: MessageEvent) => {
    const msg: WebSocketMessage = JSON.parse(event.data)
    
    if (msg.type === "edit") {
      // THE NAIVE UPDATE: Just overwrite everything
      if (msg.data !== documentText.value) {
        documentText.value = msg.data
      }
    }
  }
  
  socket.value.onerror = (error) => {
    console.error("WebSocket error:", error)
    status.value = "Disconnected"
  }
  
  socket.value.onclose = () => {
    status.value = "Disconnected"
    console.log("Disconnected from SyncPad Server")
  }
})

// Handle User Typing
const onInput = () => {
  // Send the new text to the server
  if (socket.value && socket.value.readyState === WebSocket.OPEN) {
    const message: WebSocketMessage = {
      type: "edit",
      data: documentText.value
    }
    socket.value.send(JSON.stringify(message))
  }
}

// Cleanup when closing the tab
onBeforeUnmount(() => {
  if (socket.value) {
    socket.value.close()
  }
})
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <div class="max-w-5xl mx-auto px-4 py-8 flex flex-col h-screen">
      <!-- Header -->
      <header class="mb-8">
        <div class="bg-white rounded-2xl shadow-sm border border-slate-200 px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/30">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </div>
              <div>
                <h1 class="text-2xl font-bold bg-gradient-to-r from-blue-600 to-blue-700 bg-clip-text text-transparent">
                  SyncPad
                </h1>
                <p class="text-xs text-slate-500">Collaborative Editor</p>
              </div>
            </div>
            
            <!-- Status Badge -->
            <div class="flex items-center gap-2">
              <div 
                class="px-4 py-2 rounded-full text-sm font-medium transition-all duration-300 flex items-center gap-2"
                :class="{
                  'bg-emerald-50 text-emerald-700 border border-emerald-200': status === 'Connected',
                  'bg-amber-50 text-amber-700 border border-amber-200': status === 'Connecting',
                  'bg-rose-50 text-rose-700 border border-rose-200': status === 'Disconnected'
                }"
              >
                <span 
                  class="w-2 h-2 rounded-full"
                  :class="{
                    'bg-emerald-500 animate-pulse': status === 'Connected',
                    'bg-amber-500 animate-pulse': status === 'Connecting',
                    'bg-rose-500': status === 'Disconnected'
                  }"
                ></span>
                {{ status }}
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- Editor Area -->
      <div class="flex-1 mb-6">
        <div class="h-full bg-white rounded-2xl shadow-lg border border-slate-200 overflow-hidden hover:shadow-xl transition-shadow duration-300">
          <textarea 
            v-model="documentText" 
            @input="onInput"
            placeholder="Start typing something amazing..."
            spellcheck="false"
            class="w-full h-full px-8 py-6 text-slate-800 placeholder-slate-400 resize-none outline-none focus:ring-2 focus:ring-blue-500 focus:ring-inset rounded-2xl text-base leading-relaxed transition-all"
          ></textarea>
        </div>
      </div>

      <!-- Info Footer -->
      <div class="text-center">
        <div class="inline-flex items-center gap-2 px-5 py-3 bg-white rounded-full shadow-sm border border-slate-200">
          <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-sm text-slate-600">
            Open this URL in another tab to test real-time collaboration
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom scrollbar for the textarea */
textarea::-webkit-scrollbar {
  width: 8px;
}

textarea::-webkit-scrollbar-track {
  background: transparent;
}

textarea::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}

textarea::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>