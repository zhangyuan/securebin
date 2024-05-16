<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import { httpClient } from '@/http/client';
import * as crypto from 'crypto-js';
import { useRouter } from 'vue-router'


const content = ref("");
const password = ref("");
const maxAccessCount = ref<number>();
    const messageURL = ref<string>();

const reset = () => {
    content.value = ""
    password.value = ""
    maxAccessCount.value = undefined
    messageURL.value = ""
}

const isValid = computed(() => {
  return !!content.value
})

interface CreateMessageResponse {
    id: number
}

const router = useRouter()

const onCreate = async (event: Event) => {
    event.preventDefault()

    if (content.value.length === 0) {
        return
    }

    const key = crypto.lib.WordArray.random(32).toString()

    const encrypted = crypto.AES.encrypt(content.value, key);

    const payload = {
        content: encrypted.toString(),
        password: password.value,
        max_access_count: maxAccessCount.value || 0,
    }

    const { data } = await httpClient.post<CreateMessageResponse>("/api/messages", payload)

    const messagePath = router.resolve({name: 'view', params: {id: data.id, key: key}}).href
    messageURL.value = new URL(messagePath, window.location.origin).href;
}

</script>

<template>
    <div class="container mx-auto py-4">
        <h1 class="text-center mb-2">Create Encrypted Message</h1>
        <div v-if="messageURL" class="text-center">
            <h2 class="my-4"> Please copy the URL and send it to the recipient:</h2>
            <input readonly class="text-center w-full my-1 border-solid border-2" id="message-url" v-model="messageURL" />

            <div class="my-10">
                <button @click="reset" class="border-solid border-2 bg-sky-400 text-white p-2 px-5"> Create more</button>
            </div>
        </div>
        <div v-else>
            <form>
            <div class="my-2">
                <textarea rows="10" class="form-textarea w-full" placeholder="Content" v-model="content" />
            </div>

            <div class="my-2">
                <input type="password" class="form-input w-full" readonly onfocus="this.removeAttribute('readonly');" placeholder="Password (optional)" v-model="password" />
            </div>

            <div class="my-2">
                <input type="number" class="form-input w-full" placeholder="Max access count (optional)" v-model="maxAccessCount" />
            </div>

            <div class="my-2 text-center">
                <button :disabled="!isValid" class="border-solid border-2 disabled:bg-slate-400 bg-sky-400 text-white p-2 px-5" @click="onCreate">Create</button>
            </div>
        </form>

        </div>

    </div>

</template>

<style scoped></style>
