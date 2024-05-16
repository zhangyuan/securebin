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
    password_required: boolean
}

interface ViewMessageResponse {
    content: string
    acccess_count: number
}

const message = ref<ViewMessageResponse>()
const metadata = ref<GetMessageResponse>()
const password = ref<string>()

const revealMesssage = async (password?: string) => {
    const { data: messageData } = await httpClient.post<ViewMessageResponse>(`/api/messages/${id}`, {
        password: password
    })
        message.value = messageData

        const decrypted = crypto.AES.decrypt(message.value.content, key)

        if (decrypted.sigBytes < 0) {
            error.value = "invalid key"
        }

        content.value = decrypted.toString(crypto.enc.Utf8)
}

onMounted(async () => {
    const { data: metadataData } = await httpClient.get<GetMessageResponse>(`/api/messages/${id}`)
    metadata.value = metadataData

    if (!metadata.value.password_required) {
        await revealMesssage()
    }
})



const onReveal = async (event: Event) => {
    event.preventDefault()

    if (!password) {
        return
    }

    await revealMesssage(password.value)
}


</script>

<template>
    <div class="container mx-auto py-4">
        <div v-if="error" class="text-orange-600 my-5 text-center"> {{ error }}</div>
        <div v-if="message">
            <div class="border-solid border-2 p-5"> {{ content }}</div>
            <div class="text-center my-5"> Access count: {{ message.acccess_count }}</div>
        </div>

        <div v-else-if="metadata?.password_required">
            <form>
                <div class="my-2">
                    <input type="password" class="form-input w-full" placeholder="Password" v-model="password"  />
                </div>

                <div class="my-2 text-center">
                    <button class="border-solid border-2 border-indigo-600 p-2 px-5" @click="onReveal">Reveal</button>
                </div>
            </form>
        </div>
    </div>
</template>

<style scoped></style>
