<script>
  import { pushState, searchTerms } from "../stores/url";

  let terms = $searchTerms;

  function onKeyPress({ charCode }) {
    if (charCode === 13) {
      pushState(
        { mailbox: "search", searchTerms: terms },
        `/search/?terms=${escape(terms)}`
      );
    }
  }
</script>

<div class="relative text-gray-300">
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    class="fill-current w-5 absolute left-2 top-1/2 transform -translate-y-1/2"
    ><path d="M0 0h24v24H0z" fill="none" /><path
      d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"
    /></svg
  >
  <input
    type="search"
    class="w-96 h-full outline-none pl-8 pr-1 py-1 bg-gray-600 text-sm border rounded border-gray-500 hover:border-gray-400 focus:border-gray-400"
    placeholder="Search threads"
    on:keypress={onKeyPress}
    bind:value={terms}
  />
</div>
