import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import Home from './pages/Home';
import reportWebVitals from './reportWebVitals';
import { createBrowserRouter,RouterProvider } from 'react-router-dom';
import NotFound from './components/NotFound';
import LoginForm from './components/Loginform';
import SignUpForm from './components/Signupform';
import Messenger from './components/Messenger';

const router = createBrowserRouter([
  {
    path:"/",
    element:<Home/>,
    errorElement:<NotFound/>,
    children:[
      {path:"/LogIn",element:<LoginForm/>},
      {path:"/SignUp",element:<SignUpForm/>}
    ]},
    {
      path:"/Messenger",
      element:<Messenger/>,
      errorElement: <NotFound/>
    }
])

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);

reportWebVitals();
