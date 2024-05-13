<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router'
import * as crypto from 'crypto-js';
import { httpClient } from '@/http/client';

const route = useRoute()

const id = route.params.id
const key = route.params.key.toString()
const content = ref("")

interface GetMessageResponse {
    content: string
    acccess_count: number
}

const response = ref<GetMessageResponse>()

onMounted(async () => {
    const { data } = await httpClient.post<GetMessageResponse>(`/api/messages/${id}`, {})
    response.value = data

    const decrypted = crypto.AES.decrypt(data.content, key)

    content.value = decrypted.toString(crypto.enc.Utf8)
})

</script>

<template>
        <div class="container mx-auto py-4">
            <h1 class="text-center mb-2"> Message: </h1>
            <div class="border-solid border-2 p-5"> {{ content }}</div>
            <div class="text-center my-5"> Access count: {{ response?.acccess_count }}</div>
        </div>
    
</template>

<style scoped></style>
