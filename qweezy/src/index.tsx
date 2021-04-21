import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { Global } from "@emotion/core"
import { ThemeProvider } from "emotion-theming"
import {
    defaultTheme, GlobalStyle
} from "@connctd/quartz"

import { createClient, ClientContextProvider} from 'react-fetching-library';

const client = createClient();

ReactDOM.render(
  <React.StrictMode>
    <Global styles={GlobalStyle} />
    <ThemeProvider theme={defaultTheme}>
    <ClientContextProvider client={client}>
       <App />
    </ClientContextProvider>
    </ThemeProvider>
  </React.StrictMode>,
  document.getElementById('root')
);
