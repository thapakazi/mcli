import React from "react";
import {Box, Text} from "ink";
import TextInput from "ink-text-input";

interface Props {
  value: string;
  placeholder?: string;
  onChange(input: string): void;
}

const SearchBar: React.FC<Props> = ({value, onChange, placeholder}) => (
  <Box>
    <Text>🔍 </Text>
    <TextInput
      value={value}
      placeholder={placeholder}
      onChange={onChange}
    />
  </Box>
);

export default SearchBar;
