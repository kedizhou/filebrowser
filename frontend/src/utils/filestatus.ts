import { fetchURL } from "@/api/utils";
import { baseURL } from "@/utils/constants";
export const setReadStatus = async ( name: string, path: string ) => {
    try {
      const response = await fetchURL(`/file/setfilestatus`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({name,path}),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const result = await response.status;
      if (result == 200)
        // alert("ok")
        return result
      else
        // alert("failed")
      return 409
      // return result;
    } catch (error) {
      // console.error("Error posting file status:", error);
    }
  };

  export const queryReadStatus = async ( name: string, path: string ) => {
    try {
      const response = await fetchURL(`/file/getfilestatus`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({name,path}),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const result = await response.status;
      const content = await response.text();
      if (result == 200){
        // alert("return text: "+ content)
        return content;
      }
      else
        // alert("failed")
        return ""
      
    } catch (error) {
      // console.error("Error posting file status:", error);
    }
  };

  export const queryOwner = async (  path: string ) => {
    try {
      const response = await fetchURL(`/file/getfileowner`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({path}),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const result = await response.status;
      const content = await response.text();
      if (result == 200){
        // alert("return text: "+ content)
        return content;
      }
      else
        // alert("failed")
        return ""
      
    } catch (error) {
      // console.error("Error posting file status:", error);
    }
  };