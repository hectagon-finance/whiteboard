import React, { useEffect, useState } from "react";
import elliptic from "elliptic";
import Select from "../common/Select";
import Input from "../common/Input";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { useLocation } from "react-router-dom";
import generateRandomNumber from "../function/generateRandomNumber";

const EC = elliptic.ec;
const ec = new EC("p256");

const Command = ["Create", "Start", "Stop", "Pause", "Finish", "Assign"];

const TodoList = () => {
  const schema = yup.object().shape({
    description: yup.string().required("Field is invalid"),
    title: yup.string().required("Field is invalid"),
  });
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    mode: "onChange",
    resolver: yupResolver(schema),
  });
  const location = useLocation();
  const [privateKey, setPrivateKey] = useState(null);
  const [publicKey, setPublicKey] = useState(null);
  const [walletAddress, setWalletAddress] = useState(null);

  useEffect(() => {
    if (location.state) {
      setPrivateKey(location.state.privateKey);
      setPublicKey(location.state.publicKey);
      setWalletAddress(location.state.walletAddress);
    }
  }, [location]);

  const [selectedCommand, setSelectedCommand] = useState("Create");

  console.log("privateKey: " + privateKey);
  console.log("publicKey: " + publicKey);
  console.log("walletAddress: " + walletAddress);

  const handleSubmitForm = (valueFields) => {
    let message = JSON.Stringify({
      type: "transaction",
      from: "client",
      transactionId: generateRandomNumber(),
      publicKey: publicKey,
      // signature: signatureStr,
      // data: data,
    });

    console.log(valueFields.description);
    console.log(valueFields.title);
  };

  const handleCommandChange = (newCommand) => {
    setSelectedCommand(newCommand);
    console.log("Selected command in TodoList: " + newCommand);
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Select
        listCommands={Command}
        selectedCommand={selectedCommand}
        onCommandChange={handleCommandChange}
      />
      <div className="bg-white p-10 rounded-lg shadow md:w-3/4 mx-auto lg:w-1/2 mt-[30px]">
        <form onSubmit={handleSubmit(handleSubmitForm)}>
          {selectedCommand === "Create" && (
            <div>
              <Input
                label={"Description"}
                id="description"
                register={register("description")}
                message={errors?.description?.message}
              />
              <Input
                label={"Title"}
                id="title"
                register={register("title")}
                message={errors?.title?.message}
              />
            </div>
          )}
          {/* <Input
            label={"Repo name"}
            id="repo"
            register={register("repoName")}
            message={errors?.repoName?.message}
          /> */}
          <button className="block w-full bg-blue-500 text-white font-bold p-4 rounded-lg">
            Submit
          </button>
        </form>
      </div>
    </div>
  );
};

export default TodoList;
