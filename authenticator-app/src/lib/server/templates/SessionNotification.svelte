<script lang="ts">
    import { Container, Html, Text } from 'svelte-email';
    export let userName: string = '';
    export let firstName: string = '';
    export let lastName: string = '';
    export let loginTime: string;
    export let loginLocation: string;
</script>

<Html lang="en">
    <Container>
        <Text>
            Hello {userName},
        </Text>
    </Container>
</Html>