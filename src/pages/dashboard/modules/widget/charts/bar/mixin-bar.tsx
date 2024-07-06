import React from "react";
import {
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
  { name: "Morning Shift", performance: 85, chefs: 3, servers: 10 },
  { name: "Afternoon Shift", performance: 78, chefs: 4, servers: 12 },
  { name: "Evening Shift", performance: 90, chefs: 5, servers: 8 },
];

const MixinBarChart = () => {
  return (
    <div style={{ width: "100%", height: 400 }}>
      <ResponsiveContainer>
        <BarChart
          data={data}
          margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray='3 3' />
          <XAxis dataKey='name' />
          <YAxis />
          <Tooltip />
          <Legend />
          <Bar dataKey='performance' fill='#8884d8' name='Performance' />
          <Bar dataKey='chefs' fill='#82ca9d' name='Chefs' />
          <Bar dataKey='servers' fill='#ffc658' name='Servers' />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};

export default MixinBarChart;
