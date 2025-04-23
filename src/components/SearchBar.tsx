// src/components/SearchBar.tsx
import React from "react";
import {Box, Text} from "ink";
import TextInput from "ink-text-input";

interface Props {
  value: string;
  placeholder?: string;
  onChange(input: string): void;
  focus?: boolean;
  onSubmit?(): void;      // new
}

const SearchBar: React.FC<Props> = ({
  value,
  onChange,
  placeholder,
  focus,
  onSubmit,
}) => (
  <Box>
    <Text>🔍 </Text>
    <TextInput
      value={value}
      placeholder={placeholder}
      onChange={onChange}
      focus={focus}
      onSubmit={onSubmit}    // let TextInput handle Enter
    />
  </Box>
);

export default SearchBar;
