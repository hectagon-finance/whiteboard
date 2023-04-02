import React, { useState } from "react";
import Select from "../common/Select";
import Input from "../common/Input";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

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

  const [selectedCommand, setSelectedCommand] = useState("Create");

  const handleSubmitForm = (valueFields) => {
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
