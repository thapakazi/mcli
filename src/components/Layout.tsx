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
    marginBottom={1}
    marginRight={1}
    marginLeft={1}
  >
    {header && (
      <Box>
        {header}
      </Box>
    )}

    {/* body grows to take available space if you want */}
    <Box flexDirection="column" flexGrow={1} padding={1}>
      {children}
    </Box>

    {footer && (
      <Box>
        {footer}
      </Box>
    )}
  </Box>
);

export default Layout;
