import { Suspense, use, useActionState, useState, useTransition } from "react";
import "./App.css";

type AccountData = {
  customerId: number;
  firstName: string;
  lastName: string;
  isPrimeCustomer: boolean;
};

function getAllAccountsPromise() {
  const getAllAccounts = async (): Promise<AccountData[]> => {
    try {
      const response = await fetch("http://localhost:8080/allCustomers");
      if (!response.ok) {
        throw new Error("API connection Issue");
      }

      const data = await response.json();
      return data;
    } catch (e) {
      console.log(e);
      throw e;
    }
  };

  return getAllAccounts();
}

async function addOrUpdateAccountPromise(accountData : AccountData) {
      try {
          const response = await fetch("http://localhost:8080/addCustomer", {
            method : "POST",
            headers : {
              'Content-Type' : 'application/json'
            },
            body : JSON.stringify(accountData)
          });
          
          if(!response.ok){
            throw new Error("Add Customer Failed")
          }
      } catch(e){
        console.log(e);
        throw e
      }
}

function PostList({
  getAllAccountsPromise,
  setAccountDetail,
}: {
  getAllAccountsPromise: Promise<AccountData[]>;
  setAccountDetail: (accountDetail: AccountData) => void;
}) {
  const accountList = use(getAllAccountsPromise);
  return accountList.map((item,index) => (
    <div onClick={() => setAccountDetail(accountList[index])}>{item.customerId}</div>
  ));
}

function AccountDetails({
  accountDetails
}: {
  accountDetails: AccountData | undefined;
}) {

  const [accountDetailsTemp,setAccountDetailTemp] = useState(accountDetails);
  
  const [error, submitAction, isPending] = useActionState(
    async (previousState : any,formData : FormData) => {
        try{
          console.log("formdata is :"+ formData.get("firstName"))
          await addOrUpdateAccountPromise({firstName : formData.get("firstName")?.toString()?? "",lastName:formData.get("lastName")?.toString()?? "",customerId:parseInt(formData.get("customerID")?.toString()?? "0"),isPrimeCustomer:formData.get("isPrimecustomer") === "on"})
        }catch(error){
          throw error;
        }

        return null;
    },accountDetails
  );

  if (accountDetailsTemp) {
    return (
      <form action={submitAction}>
        <input value={accountDetailsTemp.customerId} onChange={(event) => setAccountDetailTemp({...accountDetailsTemp, customerId:  parseInt(event.target.value)})} type="number" name="customerID"/>
        <input value={accountDetailsTemp.firstName} onChange={(event) => setAccountDetailTemp({...accountDetailsTemp, firstName: event.target.value})} type="text" name="firstName"/>
        <input value={accountDetailsTemp.lastName} onChange={(event) => setAccountDetailTemp({...accountDetailsTemp, lastName: event.target.value})} type="text" name="lastName"/>
        <input checked={accountDetailsTemp.isPrimeCustomer} onChange={() =>setAccountDetailTemp({...accountDetailsTemp,isPrimeCustomer:!accountDetailsTemp.isPrimeCustomer})} type="checkbox" name="isPrimecustomer"  />
        <button disabled = {isPending} type="submit"> Submit Button </button>
      </form>
    );
  } else {
    return <div>No Account Selected</div>;
  }
}


function App() {
  const [accountDetail, setAccountDetail] = useState<AccountData>();

  const updateAccountDetail = (accountDetail: AccountData) => {
    console.log("this is updating the account details")
    setAccountDetail(accountDetail);
  };
  return (
    <>
    <Suspense>
      <PostList
        getAllAccountsPromise={getAllAccountsPromise()}
        setAccountDetail={updateAccountDetail}
      />
</Suspense>
      {accountDetail ? <AccountDetails accountDetails={accountDetail}/> : <div>No Account Selected</div>}
      </>
    
  );
}

export default App;