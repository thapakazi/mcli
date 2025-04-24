// src/components/Layout.tsx
import React from "react";
import { Box } from "ink";

interface LayoutProps {
  header?: React.ReactNode;
  footer?: React.ReactNode;
  children: React.ReactNode;
}

/**
 * A simple layout with optional header, body (children), and footer.
 */
const Layout: React.FC<LayoutProps> = ({ header, children, footer }) => (
  <Box
  flexDirection="column"
  marginTop={1}
  marginRight={2} 
  marginLeft={2} 
    >
    {header && (
      <Box marginBottom={1}>
        {header}
      </Box>
    )}

    {/* body grows to take available space if you want */}
    <Box flexDirection="column" flexGrow={1}>
      {children}
    </Box>

    {footer && (
      <Box marginTop={1}>
        {footer}
      </Box>
    )}
  </Box>
);

export default Layout;
