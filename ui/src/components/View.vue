<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router'
import * as crypto from 'crypto-js';
import { httpClient } from '@/http/client';

const route = useRoute()

const id = route.params.id
const key = route.params.key.toString()
const content = ref("")
const error = ref("")

interface GetMessageResponse {
    content: string
    acccess_count: number
}

const response = ref<GetMessageResponse>()

onMounted(async () => {
    const { data } = await httpClient.post<GetMessageResponse>(`/api/messages/${id}`, {})
    response.value = data

    const decrypted = crypto.AES.decrypt(data.content, key)

    if (decrypted.sigBytes < 0) {
        error.value = "invalid key"
    }

    content.value = decrypted.toString(crypto.enc.Utf8)
})

</script>

<template>
    <div class="container mx-auto py-4">
        <div v-if="error" class="text-orange-600	my-5 text-center"> {{ error }}</div>
        <div v-else>
            <div class="border-solid border-2 p-5"> {{ content }}</div>
        </div>
        <div class="text-center my-5"> Access count: {{ response?.acccess_count }}</div>
    </div>
</template>

<style scoped></style>
