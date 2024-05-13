<script setup lang="ts">
import { RouterView } from "vue-router";
</script>

<template>
    <div>
        <div class="container mx-auto py-4">
            <ul class="flex justify-center">
                <li class="py-3 mx-2"><router-link to="/">Home</router-link></li>
            </ul>
        </div>
        <RouterView />
    </div>
</template>

<style scoped></style>
