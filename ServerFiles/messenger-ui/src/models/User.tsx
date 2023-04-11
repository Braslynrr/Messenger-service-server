class User{
    zone:string  
    number:string        
	state:string      
	userName:string     
	password:string 

    constructor(zone:string="",number:string="",state:string="",userName:string="",password:string=""){
	    this.zone=zone
        this.number=number       
	    this.state=state
	    this.userName=userName
	    this.password=password
    }
}

export default User;