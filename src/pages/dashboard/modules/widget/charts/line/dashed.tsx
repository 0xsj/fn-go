import React from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

const data = [
  {
    name: "Jan",
    male: 400,
    female: 240,
  },
  {
    name: "Feb",
    male: 300,
    female: 139,
  },
  {
    name: "Mar",
    male: 200,
    female: 980,
  },
  {
    name: "Apr",
    male: 278,
    female: 390,
  },
  {
    name: "May",
    male: 189,
    female: 480,
  },
  {
    name: "Jun",
    male: 239,
    female: 380,
  },
  {
    name: "Jul",
    male: 349,
    female: 430,
  },
];

export const DashedEmployeeCountChart = () => {
  return (
    <div>
      <ResponsiveContainer width='100%' height={350}>
        <LineChart
          data={data}
          margin={{
            top: 5,
            right: 30,
            left: 20,
            bottom: 5,
          }}
        >
          <CartesianGrid strokeDasharray='3 3' />
          <XAxis dataKey='name' />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line
            type='monotone'
            dataKey='male'
            stroke='#8884d8'
            strokeDasharray='5 5'
          />
          <Line
            type='monotone'
            dataKey='female'
            stroke='#82ca9d'
            strokeDasharray='3 4 5 2'
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};
