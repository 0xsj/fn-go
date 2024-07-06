import React from "react";
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

const data = [
  {
    month: "Jan",
    employees: 15,
    tips: 2400,
  },
  {
    month: "Feb",
    employees: 17,
    tips: 3200,
  },
  {
    month: "Mar",
    employees: 20,
    tips: 2800,
  },
  {
    month: "Apr",
    employees: 22,
    tips: 4000,
  },
  {
    month: "May",
    employees: 19,
    tips: 3600,
  },
  {
    month: "Jun",
    employees: 21,
    tips: 4500,
  },
  {
    month: "Jul",
    employees: 18,
    tips: 3800,
  },
];

export const EmployeeTipsChart = () => {
  return (
    <div>
      <ResponsiveContainer width='100%' height={400}>
        <BarChart
          data={data}
          margin={{
            top: 20,
            right: 30,
            left: 20,
            bottom: 5,
          }}
        >
          <CartesianGrid strokeDasharray='3 3' />
          <XAxis dataKey='month' />
          <YAxis yAxisId='left' orientation='left' stroke='#8884d8' />
          <YAxis yAxisId='right' orientation='right' stroke='#82ca9d' />
          <Tooltip />
          <Legend />
          <Bar yAxisId='left' dataKey='employees' fill='#8884d8' />
          <Line
            yAxisId='right'
            type='monotone'
            dataKey='tips'
            stroke='#82ca9d'
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};
