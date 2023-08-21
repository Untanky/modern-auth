<script lang="ts">
    import { Container, Html } from 'svelte-email';
    import Salutation from './components/Salutation.svelte';
    export let userName: string = '';
    export let firstName: string = '';
    export let lastName: string = '';
    export let resetLink: string;
</script>

<Html lang="en">
    <Container>
        <Salutation userName={userName} firstName={firstName} lastName={lastName} />
    </Container>
</Html>