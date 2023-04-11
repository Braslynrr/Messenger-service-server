import {useState, useEffect} from 'react';
import { Outlet,Link } from 'react-router-dom';


const Home = () => {

  const [Login,setLogin] = useState<boolean>(true)
  const [Key,setKey] = useState<string>("")

  useEffect(()=>{
    const header = new Headers()

    header.append("Content-Type","application/json")
    const body = {
       method: "GET",
       headers: header, 
    }
    fetch("/Key",body).then(promise=> promise.json()).then(data=> {
      setKey(data.initialValue)
      }
    ).catch(error => console.log(error))

  }, [])

  return (
    <div className="flex flex-wrap h-screen justify-center items-center">
    <div className="w-1/2 bg-cover flex flex-wrap">
      <div className='w-full h-full md:h-1/2'>
      <div className="px-6 py-4 w-full md:w-2/3 m-auto">
        <div className="font-bold text-xl mb-2">Messenger Service</div>
        <p className="text-base">
          Hello! We are glad to see you!
          Please log in to enjoy our services
          if your are not registered, sign up
          to be part of our comunity. 
        </p>
      </div>
      </div>
      <div className='m-auto py-10'>
      <div className="inline-flex rounded-md shadow-sm" role="group">
        <Link to="/SignUp" onClick={()=> {if(Login)setLogin(!Login)}} className={
          Login?
          "px-4 py-2 text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-l-lg hover:bg-gray-100 dark:bg-white dark:border-gray-600 dark:text-black dark:hover:text-white dark:hover:bg-gray-400 dark:focus:ring-blue-500 dark:focus:text-white":
          "px-4 py-2 text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-l-lg hover:bg-gray-100 dark:bg-gray-700 dark:border-gray-600 dark:text-white dark:hover:text-white dark:hover:bg-gray-600 dark:focus:ring-blue-500 dark:focus:text-white"
          }>
          Sign up
        </Link>
        <Link to="/LogIn" onClick={()=> {if(!Login)setLogin(!Login)}} className=
        {Login?
          "px-4 py-2 text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-r-md hover:bg-gray-100 dark:bg-gray-700 dark:border-gray-600 dark:text-white dark:hover:text-white dark:hover:bg-gray-600 dark:focus:ring-blue-500 dark:focus:text-white":
          "px-4 py-2 text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-r-md hover:bg-gray-100 dark:bg-white dark:border-gray-600 dark:text-black dark:hover:text-white dark:hover:bg-gray-400 dark:focus:ring-blue-500 dark:focus:text-white"}>
          Log in
        </Link>
      </div>
      </div>
    </div>
    <Outlet context={{Key,setKey}}/>
  </div>);  
}

export default Home;
