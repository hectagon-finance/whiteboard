const Input = ({ label, placeholder, message, register, ...props }) => {
  return (
    <div className="mb-5">
      <label className="block mb-2 font-bold text-gray-600">{label}</label>
      <input
        {...props}
        {...register}
        placeholder={placeholder}
        className="border border-gray-300 shadow p-3 w-full rounded mb-"
      />
      {message && <p className="text-sm text-red-400 mt-2">{message}</p>}
    </div>
  );
};

export default Input;
