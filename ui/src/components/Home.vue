<script setup lang="ts">
import { ref } from 'vue';
import { httpClient } from '@/http/client';
import * as crypto from 'crypto-js';
import { useRouter } from 'vue-router'


const content = ref("");
const password = ref("");
const maxAccessCount = ref<number>();

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
        password: password.value.trim(),
        max_access_count: maxAccessCount.value || 0,
    }

    const { data } = await httpClient.post<CreateMessageResponse>("/api/messages", payload)

    router.push({name: 'view', params: {id: data.id, key: key}})
}

</script>

<template>
    <div class="container mx-auto py-4">
        <h1 class="text-center mb-2">Create Encrypted Message</h1>
        <form>
            <div class="my-2">
                <textarea class="form-textarea w-full" placeholder="Content" v-model="content" />
            </div>

            <div class="my-2">
                <input type="password" class="form-input w-full" placeholder="Password (optional)" v-model="password" />
            </div>

            <div class="my-2">
                <input type="number" class="form-input w-full" placeholder="Max access count (optional)" v-model="maxAccessCount" />
            </div>

            <div class="my-2 text-center">
                <button class="border-solid border-2 border-indigo-600 p-2 px-5" @click="onCreate">Create</button>
            </div>

        </form>
    </div>

</template>

<style scoped></style>
