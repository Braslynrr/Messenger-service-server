import { SetStateAction, useEffect, useState } from "react";
import { User } from "../models/User";
import CountryCode from "../models/countrycode";
import { AiFillMinusCircle } from "react-icons/ai";
import { GiSandsOfTime } from "react-icons/gi";
import { VscPersonAdd } from "react-icons/vsc";
import Group from "../models/Group";

function NewMessageModal({ setOpenModal, user, onSetGroup }: { setOpenModal: (React.Dispatch<SetStateAction<boolean>>), user: User, onSetGroup: (group: SetStateAction<Group|undefined>) => void }) {
    const [number, setNumber] = useState<string>("")
    const [zone, setZone] = useState<string>("+506")
    const [countryCodes, setCountryCodes] = useState<CountryCode[]>([])
    const [userList, setUserList] = useState<User[]>([])
    const [error, setError] = useState<string>("")
    const [checking, setChecking] = useState<boolean>(false)

    function RemoveUser(zone: string, number: string) {
        setUserList(userList.filter(member => member.zone !== zone && member.number !== number))
    }

    function CreateGroup() {

        if (userList.length === 1) {
            onSetGroup({ id: "", admins: [], ischat: true, description: "", groupName: userList[0].username!, members: [user, ...userList] })
            setOpenModal(false)
        }
    }

    async function checkUser() {
        setChecking(true)
        if (number !== "") {
            console.log("xsd")
            if (user.zone !== zone || user.number !== number) {

                if (userList.find(member => member.zone === zone && member.number === number) === undefined) {

                    const header = new Headers()

                    header.append("Content-Type", "application/json")
                    const req = {
                        method: "GET",
                        headers: header,
                    }

                    const request = await fetch(`User/${zone}/${number}`, req)

                    const data: User = await request.json()

                    if (request.status === 200) {
                        let list = [...userList]
                        list.push(data)
                        setUserList(list)
                        setError("")
                    } else {
                        setError("Error: User is not registered!")
                    }

                } else {
                    setError(`Error: ${zone} ${number} is already added`)
                }
            } else {
                setError("Info: Your user will be added automatically")
            }

        }
        setChecking(false)
    }

    useEffect(() => {
        const header = new Headers()

        header.append("Content-Type", "application/json")
        const req = {
            method: "GET",
            headers: header,
        }

        fetch("/CountryCodes", req).then(promise => promise.json()).then(data => {
            setCountryCodes(data)

        }).catch(err => console.log(err))

    }, [])


    return (
        <div className="flex overflow-hidden shadow-lg fixed w-[20%] h-auto mt-14 ml-10 rounded-xl bg-white z-[100]">
            <div className="flex flex-col items-center m-auto space-y-2">
                {error !== "" &&
                    <div className="flex bg-yellow-400 w-full justify-center">
                        {error}
                    </div>
                }

                <div>
                    <h1 className="text-blue-700 font-semibold text-lg mb-3">Start a new chat</h1>
                </div>
                <div className="mb-4 flex space-x-2 ml-3">
                    <select value={zone} onChange={(event) => setZone(event.target.value)} className=" shadow w-1/4 border rounded py-1 px-1 bg-transparent border-opacity-50 text-black  focus:border-blue-500 focus:outline-none focus:shadow-outline">
                        {countryCodes.map(x => <option value={x.dial_code} key={x.name}>{`${x.dial_code} (${x.name})`}</option>)}
                    </select>
                    <div>
                        <label className="relative">
                            <input type='text' className="shadow appearance-none border rounded left-2 px-2 py-1 w-full bg-transparent border-opacity-50 text-black border-black focus:border-blue-500 focus:outline-none focus:shadow-outline transition duration-200"
                                value={number}
                                onChange={(event) => setNumber(event.target.value)} />
                            <span className="bg-transparent absolute -top-0.5 left-2 text-opacity-80 transition duration-200 input-text2">Number</span>
                        </label>
                    </div>
                    <div>
                        <button onClick={checkUser} disabled={checking} className=" text-white bg-blue-600 border-white rounded-md px-4 py-1 hover:bg-blue-300">{checking ? <div className="animate-pulse h-6 w-6"><GiSandsOfTime className="h-full w-full"/></div> : <div className="h-6 w-6"><VscPersonAdd className="h-full w-full" /></div>}</button>
                    </div>
                </div>
                <div>
                    <h1 className="text-blue-700 font-semibold text-lg mb-3">Users</h1>
                </div>
                <table className="table-auto items-center text-center w-[90%] px-10 border-collapse border border-slate-500">
                    <thead className=" bg-gray-400">
                        <tr>
                            <th className="border border-slate-600">Zone</th>
                            <th className="border border-slate-600">Number</th>
                            <th className="border border-slate-600">User Name</th>
                            <th className="border border-slate-600"></th>
                        </tr>
                    </thead>
                    <tbody>
                        {userList.map(user =>
                            <tr key={user.zone + user.number}>
                                <td className="border border-slate-600 ">{user.zone}</td>
                                <td className="border border-slate-600 ">{user.number}</td>
                                <td className="border border-slate-600 ">{user.username}</td>
                                <td className="border border-slate-600 text-red-800 font-semibold hover:text-red-400"><button><AiFillMinusCircle className="h-5 w-5" onClick={() => RemoveUser(user.zone, user.number)} /></button></td>
                            </tr>
                        )}
                    </tbody>
                </table>
                <div className="flex space-x-3 items-center justify-between">
                    <div className="mb-3 mt-3">
                        <button className="bg-green-600 border border-white border-opacity-50 rounded-md px-4 text-white hover:bg-green-300" type="button"
                            onClick={CreateGroup}
                        >Start</button>
                    </div>
                    <div>
                        <button className="bg-red-600 border border-white border-opacity-50 rounded-md px-4 text-white hover:bg-red-300" type="button"
                            onClick={() => setOpenModal(false)}
                        >Close</button>
                    </div>

                </div>
            </div>
        </div>
    )

}


export default NewMessageModal;
