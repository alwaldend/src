class Solution:
    def fullJustify(self, words: List[str], maxWidth: int) -> List[str]:
        words_count = len(words)
        
        def get_line(word_idx: int) -> Tuple[List[str], int]:
            line = []
            line_length = 0

            for i in range(word_idx, words_count):
                word = words[i]
                word_length = len(word)
                if line_length + word_length > maxWidth:
                    break

                line.append(word)
                line_length += word_length + 1 

            return line, line_length
        
        def create_line(line: List[str], line_length: int, word_idx: int) -> str:
            base_length = line_length - 1
            extra_spaces = maxWidth - base_length

            if len(line) == 1 or word_idx == words_count:
                return " ".join(line) + " " * extra_spaces

            word_count = len(line) - 1
            spaces_per_word = extra_spaces // word_count
            needs_extra_space = extra_spaces % word_count

            for j in range(needs_extra_space):
                line[j] += " "

            for j in range(word_count):
                line[j] += " " * spaces_per_word

            return " ".join(line)

        result = []
        i = 0
        while i < words_count:
            line, line_length = get_line(i)
            i += len(line)
            result.append(create_line(line, line_length, i))

        return result