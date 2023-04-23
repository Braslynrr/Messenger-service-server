import {useState,useEffect} from 'react';
import {User} from '../models/User'
import CountryCode from '../models/countrycode';
//import { useOutletContext } from 'react-router-dom';



function SignUpForm(){
  const [countryCodes,setCountryCodes] = useState<CountryCode[]>([])
  const [user,setUser] = useState<User>(new User("+506"))
  const [status,setStatus] = useState<number>(0)
  //const {Key} = useOutletContext<any>()

  useEffect(() => {
    const header = new Headers()

    header.append("Content-Type","application/json")
    const req = {
       method: "GET",
       headers: header, 
    }

    fetch("/CountryCodes",req).then(promise => promise.json()).then(data=>{
      setCountryCodes(data)
    }).catch(err => console.log(err))

  },[])

  function SingUp(){
    const header = new Headers()


    header.append("Content-Type","application/json")
    const body = {
       method: "POST",
       headers: header, 
       body: JSON.stringify(user)
    }

    fetch("/User",body).then(promise=>{
        setStatus(promise.status)
    }).catch(error=> console.log(error))
  }

  
   return(

        <div className='w-full md:w-1/2 h-2/3 md:h-full bg-cover flex flex-wrap bg-gray-900'>
        
        <form className="m-auto w-10/12 sm:max-w-xs p-3">
        <div className={status===0?"mb-4 w-full max-w-sm rounded-lg bg-white p-6 shadow-lg hidden":"mb-4 w-full max-w-sm rounded-lg bg-white p-6 shadow-lg block"}>

        </div>

        <div className="mb-4 flex">
          <select value={user.zone} onChange={(event)=>setUser(new User(event.target.value.split(" ")[0],user.number,user.state,user.username,user.password))} className=" shadow w-1/4 border rounded py-2 px-3 bg-gray-900 border-opacity-50 text-white  focus:border-blue-500 focus:outline-none focus:shadow-outline">
            {countryCodes.map(x=> <option value={x.dial_code} key={x.name}>{`${x.dial_code} (${x.name})`}</option>)}
          </select>
          <div className=" w-1/12"/>
          <label className='relative'>
          <input type='text' className="shadow appearance-none border rounded w-full py-2 px-3  bg-gray-900 border-opacity-50 text-white border-white focus:border-blue-500 focus:outline-none focus:shadow-outline transition duration-200"
            value={user.number}
            onChange={(event)=>setUser(new User(user.zone,event.target.value,user.state,user.username,user.password))}/>
            <span className="text-white absolute left-2 top-2 text-opacity-80 transition duration-200 input-text">Number</span>
          </label>
        </div>

        <div className="mb-4">
        <label className='relative'>
          <input type='text' className="shadow appearance-none border rounded w-full py-2 px-3  bg-gray-900 border-opacity-50 text-white border-white focus:border-blue-500 focus:outline-none focus:shadow-outline transition duration-200"
            value={user.username}
            onChange={(event)=>setUser(new User(user.zone,user.number,user.state,event.target.value,user.password))}/>
            <span className="text-white absolute left-2 -top-0.5 text-opacity-80 transition duration-200 input-text">User Name</span>
          </label>
        </div>

        <div className="mb-4">
        <label className='relative'>
          <textarea id="state" name="state" rows={3} className="shadow appearance-none border rounded w-full py-2 px-3  bg-gray-900 border-opacity-50 text-white border-white focus:border-blue-500 focus:outline-none focus:shadow-outline transition duration-200"
          value={user.state}
          onChange={(event)=>setUser(new User(user.zone,user.number,event.target.value,user.username,user.password))}/>
          <span className="text-white absolute left-2 text-opacity-80 transition duration-200 input-text">About</span>
          </label>
        </div>

        <div className="mb-4">
        <label className='relative'>
          <input type='password' className="shadow appearance-none border rounded w-full py-2 px-3  bg-gray-900 border-opacity-50 text-white border-white focus:border-blue-500 focus:outline-none focus:shadow-outline transition duration-200"
            value={user.password}
            onChange={(event)=>setUser(new User(user.zone,user.number,user.state,user.username,event.target.value))}/>
            <span className="text-white absolute left-2 -top-0.5 text-opacity-80 transition duration-200 input-text">Password</span>
          </label>

          {user.password==="" &&
            <p className="text-red-500 py-2 px-3 text-xs italic">Please input your password.</p>
          }

        </div>
        <div className="flex items-center justify-between">
          <button className="bg-white hover:bg-gray-700 text-black font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline active:bg-gray-900 w-full" type="button"
          onClick={SingUp}>
             <svg className="animate-spin h-2 w-2 mr-1 ..." viewBox="0 0 24 24">
             </svg>
             Sign In
          </button>

        </div>
        
      </form>
      </div>
    );

}

export default SignUpForm; 