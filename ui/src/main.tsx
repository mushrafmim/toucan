import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '@radix-ui/themes/styles.css'
import './index.css'
import App from './App.tsx'
import { Theme } from '@radix-ui/themes'
import { ToucanAuthProvider } from '@/shared/auth/oidc-provider'
import { RoleProvider } from '@/shared/context/role-context'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Theme
      accentColor="amber"
      appearance="light"
      grayColor="sand"
      panelBackground="solid"
      radius="large"
      scaling="100%"
    >
      <ToucanAuthProvider>
        <RoleProvider>
          <App />
        </RoleProvider>
      </ToucanAuthProvider>
    </Theme>
  </StrictMode>,
)
