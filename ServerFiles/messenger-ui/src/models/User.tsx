class User{
    zone:string  
    number:string        
	state?:string      
	username?:string     
	password?:string 
	url?:string

    constructor(zone:string="",number:string="",state:string="",username:string="",password:string=""){
	    this.zone=zone
        this.number=number       
	    this.state=state
	    this.username=username
	    this.password=password
    }

}

function equalUser(user:User,other:User):boolean{
	return other.zone===user.zone && other.number===user.number;
}

export {User,equalUser};